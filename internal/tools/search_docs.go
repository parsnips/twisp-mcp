package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/parsnips/twisp-mcp/internal/search"
)

// SearchDocsTool returns the tool definition and handler for search_docs.
func SearchDocsTool(im *search.IndexManager) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.Tool{
		Name:        "search_docs",
		Description: "Search Twisp documentation organized by Diataxis framework. Find tutorials, guides, explanations, and reference material.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search terms (e.g., 'double entry', 'ACH returns', 'velocity controls')",
				},
				"category": map[string]interface{}{
					"type":        "string",
					"description": "Filter by documentation category",
					"enum":        []string{"introduction", "accounting-core", "tutorials", "guides", "reference", "infrastructure", "processors"},
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum number of results (default 10)",
				},
			},
			Required: []string{"query"},
		},
	}

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleSearchDocs(ctx, request, im)
	}

	return tool, handler
}

func handleSearchDocs(ctx context.Context, request mcp.CallToolRequest, im *search.IndexManager) (*mcp.CallToolResult, error) {
	query, ok := request.Params.Arguments["query"].(string)
	if !ok || query == "" {
		return mcp.NewToolResultError("query is required"), nil
	}

	var category string
	if c, ok := request.Params.Arguments["category"].(string); ok {
		category = c
	}

	limit := 10
	if l, ok := request.Params.Arguments["limit"].(float64); ok {
		limit = int(l)
	}

	docs, err := im.SearchDocs(query, category, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("search failed: %v", err)), nil
	}

	if len(docs) == 0 {
		return mcp.NewToolResultText("No matching documentation found."), nil
	}

	// Format results - truncate content for summary
	results := make([]map[string]interface{}, len(docs))
	for i, doc := range docs {
		excerpt := doc.Content
		if len(excerpt) > 500 {
			excerpt = excerpt[:500] + "..."
		}
		results[i] = map[string]interface{}{
			"id":          doc.ID,
			"title":       doc.Title,
			"description": doc.Description,
			"path":        doc.Path,
			"category":    doc.Category,
			"excerpt":     excerpt,
		}
	}

	output, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to format results: %v", err)), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

// GetDocTool returns the tool definition and handler for get_doc.
func GetDocTool(im *search.IndexManager) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.Tool{
		Name:        "get_doc",
		Description: "Retrieve a full documentation page by its path.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "Documentation path (e.g., 'accounting-core/balances', 'guides/velocity-controls')",
				},
			},
			Required: []string{"path"},
		},
	}

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetDoc(ctx, request, im)
	}

	return tool, handler
}

func handleGetDoc(ctx context.Context, request mcp.CallToolRequest, im *search.IndexManager) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok || path == "" {
		return mcp.NewToolResultError("path is required"), nil
	}

	doc, err := im.GetDocDocument(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("document not found: %v", err)), nil
	}

	result := map[string]interface{}{
		"title":       doc.Title,
		"description": doc.Description,
		"path":        doc.Path,
		"category":    doc.Category,
		"content":     doc.Content,
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to format result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}
