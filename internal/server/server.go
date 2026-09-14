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

const Version = "3.1.0"

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
	cloudOnly bool
}

func New(local *graphql.Client, upstream Upstream) *Server {
	s := &Server{upstream: upstream, cloudOnly: local == nil}
	if local != nil {
		tool, handler := tools.GraphQLExecuteTool(local)
		s.local = server.ServerTool{Tool: tool, Handler: handler}
	}
	hooks := &server.Hooks{}
	hooks.AddBeforeListTools(func(ctx context.Context, _ any, _ *mcp.ListToolsRequest) { s.refresh(ctx) })
	s.mcp = server.NewMCPServer("twisp-mcp", Version, server.WithHooks(hooks), server.WithToolCapabilities(true))
	if !s.cloudOnly {
		s.mcp.AddTools(s.local)
	}
	return s
}

// Refresh on tools/list, retaining the last successful catalog during an
// outage. GraphQL registration never depends on cloud discovery succeeding.
func (s *Server) refresh(ctx context.Context) { _ = s.refreshWithin(ctx, 5*time.Second) }

func (s *Server) refreshWithin(ctx context.Context, timeout time.Duration) error {
	s.refreshMu.Lock()
	defer s.refreshMu.Unlock()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	result, err := s.upstream.ListTools(ctx)
	if err != nil {
		log.Print("cloud tool discovery unavailable; retaining existing tools")
		return err
	}
	catalog := make([]mcp.Tool, 0, len(result.Tools))
	for _, tool := range result.Tools {
		if s.cloudOnly || tool.Name != "graphql_execute" {
			catalog = append(catalog, tool)
		}
	}
	if reflect.DeepEqual(catalog, s.catalog) {
		return nil
	}
	s.catalog = catalog
	registered := []server.ServerTool{}
	if !s.cloudOnly {
		registered = append(registered, s.local)
	}
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
	return nil
}

func (s *Server) Serve() error { return server.ServeStdio(s.mcp) }

// CheckCloud establishes authenticated access and fills the initial catalog.
// Pure cloud startup must not silently succeed with no usable tools.
func (s *Server) CheckCloud(ctx context.Context) error {
	return s.refreshWithin(ctx, 30*time.Second)
}
