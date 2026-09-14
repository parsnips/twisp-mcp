package cloud

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestIAMProofProtocol(t *testing.T) {
	now := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)
	issuer := "https://auth.us-east-1.dev.twisp.com/"
	creds := aws.Credentials{AccessKeyID: "AKIDEXAMPLE", SecretAccessKey: "example-secret", SessionToken: "session-token"}
	proof, err := iamProof(context.Background(), creds, "us-east-1", issuer, now)
	if err != nil {
		t.Fatal(err)
	}
	var envelope struct {
		Token      string
		Expiration time.Time
	}
	if err := json.Unmarshal(proof, &envelope); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(envelope.Token, "twisp-aws-v1.") || !envelope.Expiration.Equal(now.Add(14*time.Minute)) {
		t.Fatal("invalid proof envelope")
	}
	raw, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(envelope.Token, "twisp-aws-v1."))
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(string(raw))
	if err != nil {
		t.Fatal(err)
	}
	if u.Scheme != "https" || u.Host != "sts.us-east-1.amazonaws.com" || u.Path != "/" {
		t.Fatal("invalid STS URL")
	}
	expected := map[string]string{
		"Action": "GetCallerIdentity", "Version": "2011-06-15", "X-Amz-Expires": "60",
		"X-Amz-Algorithm": "AWS4-HMAC-SHA256", "X-Amz-Credential": "AKIDEXAMPLE/20260914/us-east-1/sts/aws4_request",
		"X-Amz-Date": "20260914T120000Z", "X-Amz-Security-Token": "session-token",
		"X-Amz-SignedHeaders": "host;x-twisp-aws-id",
	}
	q := u.Query()
	for key, value := range expected {
		if q.Get(key) != value {
			t.Errorf("%s = %q, want %q", key, q.Get(key), value)
		}
	}
	if len(q) != len(expected)+1 || len(q.Get("X-Amz-Signature")) != 64 {
		t.Fatal("unexpected signing parameters")
	}
	// A different issuer must produce a different proof, even with the same AWS identity.
	other, err := iamProof(context.Background(), creds, "us-east-1", "https://auth.us-east-1.cloud.twisp.com/", now)
	if err != nil || string(other) == string(proof) {
		t.Fatal("proof not bound to issuer")
	}
}

func TestTokenRefreshAndFailure(t *testing.T) {
	now := time.Now()
	var exchanges atomic.Int32
	source := &tokenSource{now: func() time.Time { return now }}
	source.exchange = func(context.Context) (string, time.Time, error) {
		n := exchanges.Add(1)
		if n >= 3 {
			return "", time.Time{}, fmt.Errorf("expired AWS session")
		}
		return fmt.Sprintf("jwt-%d", n), now.Add(time.Hour), nil
	}
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token, err := source.Token(context.Background())
			if err != nil || token != "jwt-1" {
				t.Error(token, err)
			}
		}()
	}
	wg.Wait()
	if exchanges.Load() != 1 {
		t.Fatal("concurrent requests repeated exchange")
	}
	now = now.Add(59 * time.Minute)
	if token, err := source.Token(context.Background()); err != nil || token != "jwt-2" {
		t.Fatal(token, err)
	}
	now = now.Add(time.Hour)
	if token, err := source.Token(context.Background()); err == nil || token != "" {
		t.Fatal("used stale token after renewal failed")
	}
}

func TestExplicitTokenPrecedence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "token")
	source := &tokenSource{config: Config{BearerToken: "Bearer env-token", TokenFile: path}, exchange: func(context.Context) (string, time.Time, error) {
		t.Fatal("fell back to AWS identity")
		return "", time.Time{}, nil
	}}
	if _, err := source.Token(context.Background()); err == nil {
		t.Fatal("missing file fell back")
	}
	for _, value := range []string{"file-token", "Bearer refreshed"} {
		if err := os.WriteFile(path, []byte(value), 0600); err != nil {
			t.Fatal(err)
		}
		actual, err := source.Token(context.Background())
		if err != nil || actual != strings.TrimPrefix(value, "Bearer ") {
			t.Fatal(actual, err)
		}
	}
	if err := os.WriteFile(path, nil, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := source.Token(context.Background()); err == nil {
		t.Fatal("empty file fell back")
	}
	source.config.TokenFile = ""
	if token, err := source.Token(context.Background()); err != nil || token != "env-token" {
		t.Fatal(token, err)
	}
}

func TestExchangeResponses(t *testing.T) {
	for _, test := range []struct {
		name, body string
		status     int
		valid      bool
	}{
		{"valid", "header." + base64.RawURLEncoding.EncodeToString([]byte(`{"exp":2000000000}`)) + ".signature", 200, true},
		{"missing-exp", "header.e30.signature", 200, false},
		{"invalid", "secret-invalid-body", 200, false},
		{"rejected", "secret-error-body", 403, false},
		{"oversize", strings.Repeat("secret", 12000), 200, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			host := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" || r.Header.Get("Content-Type") != "application/json" {
					t.Error("invalid exchange request")
				}
				w.WriteHeader(test.status)
				fmt.Fprint(w, test.body)
			}))
			defer host.Close()
			value, expiry, err := exchangeProof(context.Background(), host.Client(), host.URL, []byte(`{"Token":"proof"}`))
			if test.valid {
				if err != nil || value != test.body || expiry.Unix() != 2000000000 {
					t.Fatal(value, expiry, err)
				}
			} else if err == nil || strings.Contains(err.Error(), "secret") {
				t.Fatal("invalid response accepted or body leaked", err)
			}
		})
	}
}

func TestIAMDoesNotMintForCustomURL(t *testing.T) {
	s := &tokenSource{config: Config{URL: "https://example.com/mcp", Environment: "dev", Region: "us-east-1"}, now: time.Now}
	if _, err := s.Token(context.Background()); err == nil || !strings.Contains(err.Error(), "custom MCP URLs") {
		t.Fatal(err)
	}
}
