package cloud

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
)

// tokenSource keeps renewable credentials and JWTs in memory. Explicit tokens
// retain precedence and never silently fall back to another identity.
type tokenSource struct {
	config      Config
	mu          sync.Mutex
	value       string
	expires     time.Time
	credentials aws.CredentialsProvider
	now         func() time.Time
	exchange    func(context.Context) (string, time.Time, error)
}

func (s *tokenSource) Token(ctx context.Context) (string, error) {
	if s.config.TokenFile != "" {
		data, err := os.ReadFile(s.config.TokenFile)
		if err != nil {
			return "", fmt.Errorf("cannot read TWISP_MCP_TOKEN_FILE")
		}
		return cleanToken(string(data))
	}
	if s.config.BearerToken != "" {
		return cleanToken(s.config.BearerToken)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.value != "" && s.now().Add(time.Minute).Before(s.expires) {
		return s.value, nil
	}
	exchange := s.exchange
	if exchange == nil {
		exchange = s.exchangeIAM
	}
	value, expires, err := exchange(ctx)
	if err != nil {
		return "", err
	}
	if !s.now().Add(time.Minute).Before(expires) {
		return "", fmt.Errorf("Twisp returned a token that expires too soon")
	}
	s.value, s.expires = value, expires
	return value, nil
}

func cleanToken(value string) (string, error) {
	value = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(value), "Bearer "))
	if value == "" {
		return "", fmt.Errorf("configured cloud token is empty")
	}
	return value, nil
}

var environmentName = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
var regionName = regexp.MustCompile(`^[a-z]{2}-[a-z]+-[0-9]+$`)

func (s *tokenSource) exchangeIAM(ctx context.Context) (string, time.Time, error) {
	c := s.config
	if !environmentName.MatchString(c.Environment) || !regionName.MatchString(c.Region) {
		return "", time.Time{}, fmt.Errorf("AWS token exchange requires a valid --env and --region in the commercial AWS partition")
	}
	// An ambient IAM identity may only be forwarded to its selected Twisp target.
	if c.URL != Endpoint(c.Environment, c.Region) {
		return "", time.Time{}, fmt.Errorf("custom MCP URLs require an explicit cloud token or token file")
	}
	if s.credentials == nil {
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(c.Region))
		if err != nil {
			return "", time.Time{}, fmt.Errorf("cannot load AWS configuration; check your AWS profile")
		}
		s.credentials = cfg.Credentials
	}
	credentials, err := s.credentials.Retrieve(ctx)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("cannot obtain AWS credentials; sign in to your AWS profile or renew your AWS session")
	}
	now := s.now()
	if credentials.CanExpire && !now.Before(credentials.Expires) {
		return "", time.Time{}, fmt.Errorf("AWS session expired; renew it and relaunch the agent")
	}
	issuer := fmt.Sprintf("https://auth.%s.%s.twisp.com/", c.Region, c.Environment)
	proof, err := iamProof(ctx, credentials, c.Region, issuer, now)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("cannot sign AWS identity proof")
	}
	return exchangeProof(ctx, &http.Client{Timeout: 30 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}, issuer+"token/iam", proof)
}

// This is Twisp's twisp-aws-v1 protocol: a presigned GetCallerIdentity URL
// bound to the selected issuer with the signed x-twisp-aws-id header.
func iamProof(ctx context.Context, credentials aws.Credentials, region, issuer string, now time.Time) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://sts."+region+".amazonaws.com/?Action=GetCallerIdentity&Version=2011-06-15&X-Amz-Expires=60", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-twisp-aws-id", issuer)
	hash := sha256.Sum256(nil)
	signed, _, err := v4.NewSigner().PresignHTTP(ctx, credentials, req, fmt.Sprintf("%x", hash), "sts", region, now)
	if err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		Token      string
		Expiration time.Time
	}{
		Token:      "twisp-aws-v1." + base64.RawURLEncoding.EncodeToString([]byte(signed)),
		Expiration: now.Add(14 * time.Minute),
	})
}

func exchangeProof(ctx context.Context, client *http.Client, endpoint string, proof []byte) (string, time.Time, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(proof))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid token exchange endpoint")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("Twisp token exchange failed; check connectivity")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("Twisp token exchange rejected AWS identity (HTTP %d)", resp.StatusCode)
	}
	const maxToken = 64 << 10
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxToken+1))
	if err != nil || len(data) > maxToken {
		return "", time.Time{}, fmt.Errorf("invalid token exchange response")
	}
	value, err := cleanToken(string(data))
	if err != nil {
		return "", time.Time{}, err
	}
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return "", time.Time{}, fmt.Errorf("Twisp returned an invalid JWT")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", time.Time{}, fmt.Errorf("Twisp returned an invalid JWT")
	}
	var claims struct {
		Exp int64 `json:"exp"`
	}
	// Read exp only for refresh scheduling; the cloud API verifies the JWT.
	if json.Unmarshal(payload, &claims) != nil || claims.Exp <= 0 {
		return "", time.Time{}, fmt.Errorf("Twisp returned a JWT without a valid expiry")
	}
	return value, time.Unix(claims.Exp, 0), nil
}
