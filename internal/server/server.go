package server

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/parsnips/twisp-mcp/internal/search"
	"github.com/parsnips/twisp-mcp/internal/tools"
)

// Server wraps the MCP server with Twisp-specific functionality.
type Server struct {
	mcp *server.MCPServer
}

// New creates a new Twisp MCP server.
func New(im *search.IndexManager) (*Server, error) {
	mcpServer := server.NewMCPServer(
		"twisp-mcp",
		"2.0.0",
	)

	// GraphQL execute tool
	tool, handler := tools.GraphQLExecuteTool()
	mcpServer.AddTool(tool, handler)

	// Search tools
	schemaTool, schemaHandler := tools.SearchSchemaTool(im)
	mcpServer.AddTool(schemaTool, schemaHandler)

	docsTool, docsHandler := tools.SearchDocsTool(im)
	mcpServer.AddTool(docsTool, docsHandler)

	getDocTool, getDocHandler := tools.GetDocTool(im)
	mcpServer.AddTool(getDocTool, getDocHandler)

	examplesTool, examplesHandler := tools.SearchExamplesTool(im)
	mcpServer.AddTool(examplesTool, examplesHandler)

	getExampleTool, getExampleHandler := tools.GetExampleTool(im)
	mcpServer.AddTool(getExampleTool, getExampleHandler)

	return &Server{mcp: mcpServer}, nil
}

// Serve starts the MCP server on stdio.
func (s *Server) Serve() error {
	return server.ServeStdio(s.mcp)
}
