package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/parsnips/twisp-mcp/internal/search"
)

// SearchSchemaTool returns the tool definition and handler for search_schema.
func SearchSchemaTool(im *search.IndexManager) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.Tool{
		Name:        "search_schema",
		Description: "Search the Twisp GraphQL schema for types, mutations, queries, and fields. Use this to discover available API operations and data structures.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search terms (e.g., 'create account', 'balance', 'transaction')",
				},
				"kind": map[string]interface{}{
					"type":        "string",
					"description": "Filter by schema element kind",
					"enum":        []string{"type", "input", "enum", "mutation", "query", "scalar", "directive", "interface", "union"},
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
		return handleSearchSchema(ctx, request, im)
	}

	return tool, handler
}

func handleSearchSchema(ctx context.Context, request mcp.CallToolRequest, im *search.IndexManager) (*mcp.CallToolResult, error) {
	query, ok := request.Params.Arguments["query"].(string)
	if !ok || query == "" {
		return mcp.NewToolResultError("query is required"), nil
	}

	var kind string
	if k, ok := request.Params.Arguments["kind"].(string); ok {
		kind = k
	}

	limit := 10
	if l, ok := request.Params.Arguments["limit"].(float64); ok {
		limit = int(l)
	}

	docs, err := im.SearchSchema(query, kind, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("search failed: %v", err)), nil
	}

	if len(docs) == 0 {
		return mcp.NewToolResultText("No matching schema elements found."), nil
	}

	// Format results
	results := make([]map[string]interface{}, len(docs))
	for i, doc := range docs {
		results[i] = map[string]interface{}{
			"id":          doc.ID,
			"kind":        doc.Kind,
			"name":        doc.Name,
			"description": doc.Description,
			"sourceFile":  doc.SourceFile,
			"definition":  doc.Definition,
		}
	}

	output, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to format results: %v", err)), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}
