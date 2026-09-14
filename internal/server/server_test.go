package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	sdk "github.com/mark3labs/mcp-go/server"
	"github.com/parsnips/twisp-mcp/internal/cloud"
	"github.com/parsnips/twisp-mcp/internal/graphql"
)

func TestLocalGraphQLAndCloudTools(t *testing.T) {
	ctx := context.Background()
	var cloudCalls, localCalls atomic.Int32
	var offline atomic.Bool
	var expectedToken atomic.Value
	expectedToken.Store("Bearer cloud-token")
	remote := sdk.NewMCPServer("cloud", "test")
	remote.AddTool(mcp.NewTool("graphql_execute"), func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		t.Error("GraphQL reached cloud")
		return mcp.NewToolResultText("wrong"), nil
	})
	expected := mcp.NewToolResultText("cloud result")
	expected.Meta = &mcp.Meta{AdditionalFields: map[string]any{"twisp/embeddingBytes": float64(5)}}
	remote.AddTool(mcp.NewTool("search_docs", mcp.WithDescription("Cloud search"), mcp.WithString("query", mcp.Required())), func(_ context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cloudCalls.Add(1)
		if r.GetString("query", "") != "hello" {
			t.Error("arguments changed")
		}
		return expected, nil
	})
	transport := sdk.NewStreamableHTTPServer(remote, sdk.WithStateLess(true))
	host := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != expectedToken.Load().(string) || r.Header.Get("X-Twisp-Account-Id") != "cloud-account" {
			t.Error("wrong cloud credentials")
		}
		if offline.Load() {
			http.Error(w, "offline", 503)
			return
		}
		transport.ServeHTTP(w, r)
	}))
	defer host.Close()
	local := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		localCalls.Add(1)
		if r.Header.Get("Authorization") != "Bearer local-token" || r.Header.Get("X-Twisp-Account-Id") != "local-account" {
			t.Error("cloud credentials leaked to GraphQL")
		}
		var body graphql.GraphQLRequest
		if json.NewDecoder(r.Body).Decode(&body) != nil || body.Query != "query Local { hello }" || body.OperationName != "Local" || body.Variables["value"] != float64(7) {
			t.Error("GraphQL arguments changed")
		}
		_, _ = w.Write([]byte(`{"data":{"hello":"local"}}`))
	}))
	defer local.Close()
	tokenFile := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(tokenFile, []byte("cloud-token\n"), 0600); err != nil {
		t.Fatal(err)
	}
	upstream, err := cloud.New(cloud.Config{URL: host.URL, AccountID: "cloud-account", TokenFile: tokenFile})
	if err != nil {
		t.Fatal(err)
	}
	defer upstream.Close()
	s := New(graphql.NewClientWithConfig(local.URL, "local-account", "", "local-token"), upstream)
	c, err := client.NewInProcessClient(s.mcp)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	init := mcp.InitializeRequest{}
	init.Params.ProtocolVersion = "2025-03-26"
	if _, err = c.Initialize(ctx, init); err != nil {
		t.Fatal(err)
	}
	list, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil || len(list.Tools) != 2 {
		t.Fatal(list, err)
	}
	call := mcp.CallToolRequest{}
	call.Params.Name = "search_docs"
	call.Params.Arguments = map[string]any{"query": "hello"}
	result, err := c.CallTool(ctx, call)
	if err != nil || !reflect.DeepEqual(result, expected) {
		t.Fatalf("%+v %v", result, err)
	}
	// Token files are re-read without restarting the local process.
	expectedToken.Store("Bearer rotated")
	if err := os.WriteFile(tokenFile, []byte("Bearer rotated"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = c.CallTool(ctx, call); err != nil {
		t.Fatal(err)
	}
	offline.Store(true)
	result, err = c.CallTool(ctx, call)
	if err != nil || !result.IsError {
		t.Fatal("cloud failure not a tool error", result, err)
	}
	if list, err = c.ListTools(ctx, mcp.ListToolsRequest{}); err != nil || len(list.Tools) != 2 {
		t.Fatal("cached catalog lost", err)
	}
	call.Params.Name = "graphql_execute"
	call.Params.Arguments = map[string]any{"query": "query Local { hello }", "operationName": "Local", "variables": map[string]any{"value": 7}}
	result, err = c.CallTool(ctx, call)
	if err != nil || result.IsError || localCalls.Load() != 1 || cloudCalls.Load() != 2 {
		t.Fatal("local routing failed", result, err)
	}
	// A fresh local server must also work when initial discovery fails.
	fresh := New(graphql.NewClientWithConfig(local.URL, "local-account", "", "local-token"), upstream)
	fresh.refresh(ctx)
	if len(fresh.mcp.ListTools()) != 1 {
		t.Fatal("offline startup should expose local GraphQL")
	}
	offline.Store(false)
	remote.AddTool(mcp.NewTool("new_cloud_tool"), func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultError("upstream tool error"), nil
	})
	remote.DeleteTools("search_docs")
	list, err = c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil || len(list.Tools) != 2 {
		t.Fatal("discovery did not recover", list, err)
	}
	call.Params.Name = "new_cloud_tool"
	call.Params.Arguments = map[string]any{}
	result, err = c.CallTool(ctx, call)
	if err != nil || !result.IsError {
		t.Fatal("upstream tool error lost", err)
	}
	if s.mcp.GetTool("search_docs") != nil {
		t.Fatal("removed cloud tool retained")
	}
}
