package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestSSEAndCredentialIsolation(t *testing.T) {
	var calls atomic.Int32
	host := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer cloud" || r.Header.Get("X-Twisp-Account-Id") != "tenant" {
			t.Error("wrong cloud identity")
		}
		var request struct {
			ID     json.RawMessage `json:"id"`
			Method string          `json:"method"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if len(request.ID) == 0 {
			w.WriteHeader(202)
			return
		}
		var result string
		switch request.Method {
		case "initialize":
			result = `{"protocolVersion":"2025-03-26","capabilities":{"tools":{}},"serverInfo":{"name":"cloud","version":"1"}}`
		case "tools/call":
			calls.Add(1)
			result = `{"content":[{"type":"text","text":"answer"}],"_meta":{"twisp/embeddingBytes":5}}`
		default:
			t.Errorf("unexpected method %s", request.Method)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprintf(w, "event: message\ndata: {\"jsonrpc\":\"2.0\",\"id\":%s,\"result\":%s}\n\n", request.ID, result)
	}))
	defer host.Close()
	c, err := New(Config{URL: host.URL, AccountID: "tenant", BearerToken: "cloud"})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	request := mcp.CallToolRequest{}
	request.Params.Name = "search_docs"
	request.Params.Arguments = map[string]any{"query": "hello"}
	request.Header = http.Header{"Authorization": {"Bearer local-secret"}, "X-Twisp-Account-Id": {"wrong"}}
	result, err := c.CallTool(context.Background(), request)
	if err != nil || result.Meta.AdditionalFields["twisp/embeddingBytes"] != float64(5) {
		t.Fatal(result, err)
	}
	request.Params.Name = "graphql_execute"
	if _, err = c.CallTool(context.Background(), request); err == nil || calls.Load() != 1 {
		t.Fatal("GraphQL was forwarded")
	}
}

func TestCloudRedirectDoesNotForwardCredentials(t *testing.T) {
	var reached atomic.Bool
	target := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { reached.Store(true) }))
	defer target.Close()
	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, target.URL, 307) }))
	defer redirect.Close()
	c, err := New(Config{URL: redirect.URL, AccountID: "tenant", BearerToken: "secret"})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if _, err := c.ListTools(context.Background()); err == nil || reached.Load() {
		t.Fatal("followed credential redirect")
	}
}

func TestCanceledCallDoesNotWaitForBusyCloud(t *testing.T) {
	c, err := New(Config{URL: DefaultURL})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.acquire(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer c.release()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if _, err := c.ListTools(ctx); err == nil {
		t.Fatal("ignored cancellation")
	}
}

func TestCloudURLValidation(t *testing.T) {
	for _, endpoint := range []string{"http://example.com/mcp", "https://user:secret@example.com/mcp", ":bad"} {
		if _, err := New(Config{URL: endpoint}); err == nil {
			t.Fatal("accepted", endpoint)
		}
	}
	for _, endpoint := range []string{DefaultURL, "http://127.0.0.1:8000/mcp", "http://[::1]:8000/mcp"} {
		if _, err := New(Config{URL: endpoint}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestHTTPDiagnosticsDoNotExposeResponseBodies(t *testing.T) {
	for _, test := range []struct {
		status     int
		body, hint string
	}{
		{401, "unauthorized", "load balancer"},
		{401, "secret-reflected-token", "token expiry"},
		{400, "secret-reflected-token", "alias/"},
		{403, "secret-reflected-token", "Twisp client registration"},
		{404, "secret-reflected-token", "deployed"},
		{502, "secret-reflected-token", "runtime is unavailable"},
		{307, "secret-reflected-token", "endpoint"},
	} {
		t.Run(fmt.Sprint(test.status)+test.hint, func(t *testing.T) {
			host := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(test.status); fmt.Fprint(w, test.body) }))
			defer host.Close()
			c, err := New(Config{URL: host.URL, AccountID: "tenant", BearerToken: "secret-reflected-token"})
			if err != nil {
				t.Fatal(err)
			}
			defer c.Close()
			_, err = c.ListTools(context.Background())
			if err == nil || !strings.Contains(err.Error(), fmt.Sprintf("HTTP %d", test.status)) || !strings.Contains(err.Error(), test.hint) || strings.Contains(err.Error(), "secret-reflected-token") {
				t.Fatalf("wrong diagnostic: %v", err)
			}
		})
	}
}
