// Package server serves local MCP over stdio and routes tools by name.
package server

import (
	"context"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/parsnips/twisp-mcp/internal/graphql"
	"github.com/parsnips/twisp-mcp/internal/tools"
)

const Version = "3.0.0"

type Upstream interface {
	ListTools(context.Context) (*mcp.ListToolsResult, error)
	CallTool(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

type Server struct {
	mcp       *server.MCPServer
	upstream  Upstream
	local     server.ServerTool
	refreshMu sync.Mutex
	catalog   []mcp.Tool
}

func New(local *graphql.Client, upstream Upstream) *Server {
	tool, handler := tools.GraphQLExecuteTool(local)
	s := &Server{upstream: upstream, local: server.ServerTool{Tool: tool, Handler: handler}}
	hooks := &server.Hooks{}
	hooks.AddBeforeListTools(func(ctx context.Context, _ any, _ *mcp.ListToolsRequest) { s.refresh(ctx) })
	s.mcp = server.NewMCPServer("twisp-mcp", Version, server.WithHooks(hooks), server.WithToolCapabilities(true))
	s.mcp.AddTools(s.local)
	return s
}

// Refresh on tools/list, retaining the last successful catalog during an
// outage. GraphQL registration never depends on cloud discovery succeeding.
func (s *Server) refresh(ctx context.Context) {
	s.refreshMu.Lock()
	defer s.refreshMu.Unlock()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result, err := s.upstream.ListTools(ctx)
	if err != nil {
		log.Print("cloud tool discovery unavailable; local GraphQL remains available")
		return
	}
	catalog := make([]mcp.Tool, 0, len(result.Tools))
	for _, tool := range result.Tools {
		if tool.Name != "graphql_execute" {
			catalog = append(catalog, tool)
		}
	}
	if reflect.DeepEqual(catalog, s.catalog) {
		return
	}
	s.catalog = catalog
	registered := []server.ServerTool{s.local}
	for _, tool := range catalog {
		name := tool.Name
		registered = append(registered, server.ServerTool{Tool: tool, Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
			defer cancel()
			request.Params.Name = name
			result, err := s.upstream.CallTool(ctx, request)
			if err != nil {
				return mcp.NewToolResultError("Cloud tool unavailable. Check cloud connectivity and TWISP_MCP credentials."), nil
			}
			return result, nil
		}})
	}
	s.mcp.SetTools(registered...)
}

func (s *Server) Serve() error { return server.ServeStdio(s.mcp) }
