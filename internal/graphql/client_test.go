package graphql

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestLocalDefaultsIgnoreCloudCredentials(t *testing.T) {
	t.Setenv("TWISP_ENDPOINT", "")
	t.Setenv("TWISP_ACCOUNT_ID", "")
	t.Setenv("TWISP_BEARER_TOKEN", "")
	t.Setenv("TWISP_API_KEY", "")
	t.Setenv("TWISP_MCP_BEARER_TOKEN", "cloud-secret")
	t.Setenv("TWISP_MCP_ACCOUNT_ID", "cloud-account")
	c := NewClient()
	if c.endpoint != "http://localhost:8080/financial/v1/graphql" || c.accountID != "000000000000" || c.bearerToken != "" || c.apiKey != "" {
		t.Fatal("cloud configuration leaked to GraphQL")
	}
}
func TestGraphQLDoesNotFollowRedirects(t *testing.T) {
	var reached atomic.Bool
	target := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { reached.Store(true) }))
	defer target.Close()
	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, target.URL, 307) }))
	defer redirect.Close()
	c := NewClientWithConfig(redirect.URL, "local", "", "local-secret")
	if _, err := c.Execute(context.Background(), "query { hello }", nil, ""); err == nil || reached.Load() {
		t.Fatal("followed GraphQL redirect")
	}
}
