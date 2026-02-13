package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/parsnips/twisp-mcp/internal/search"
)

// SearchExamplesTool returns the tool definition and handler for search_examples.
func SearchExamplesTool(im *search.IndexManager) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.Tool{
		Name:        "search_examples",
		Description: "Search 585 GraphQL example documents. Find working code examples for common operations.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search terms (e.g., 'post transaction', 'multi-currency', 'create account')",
				},
				"category": map[string]interface{}{
					"type":        "string",
					"description": "Filter by example category",
					"enum":        []string{"examples", "reference", "fixtures", "tranCodeLibrary"},
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum number of results (default 5)",
				},
			},
			Required: []string{"query"},
		},
	}

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleSearchExamples(ctx, request, im)
	}

	return tool, handler
}

func handleSearchExamples(ctx context.Context, request mcp.CallToolRequest, im *search.IndexManager) (*mcp.CallToolResult, error) {
	query, ok := request.Params.Arguments["query"].(string)
	if !ok || query == "" {
		return mcp.NewToolResultError("query is required"), nil
	}

	var category string
	if c, ok := request.Params.Arguments["category"].(string); ok {
		category = c
	}

	limit := 5
	if l, ok := request.Params.Arguments["limit"].(float64); ok {
		limit = int(l)
	}

	docs, err := im.SearchExamples(query, category, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("search failed: %v", err)), nil
	}

	if len(docs) == 0 {
		return mcp.NewToolResultText("No matching examples found."), nil
	}

	// Format results - include the query
	results := make([]map[string]interface{}, len(docs))
	for i, doc := range docs {
		results[i] = map[string]interface{}{
			"id":          doc.ID,
			"category":    doc.Category,
			"subcategory": doc.Subcategory,
			"name":        doc.Name,
			"query":       doc.Query,
		}
		if doc.Variables != "" {
			results[i]["variables"] = doc.Variables
		}
		if doc.Instructions != "" {
			results[i]["instructions"] = doc.Instructions
		}
	}

	output, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to format results: %v", err)), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

// GetExampleTool returns the tool definition and handler for get_example.
func GetExampleTool(im *search.IndexManager) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.Tool{
		Name:        "get_example",
		Description: "Retrieve a complete example by its path, including request.gql, variables.json, and response.json.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "Example path (e.g., 'examples/basicTranCodeFlow/001_BasicTranCodeSetup')",
				},
			},
			Required: []string{"path"},
		},
	}

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetExample(ctx, request, im)
	}

	return tool, handler
}

func handleGetExample(ctx context.Context, request mcp.CallToolRequest, im *search.IndexManager) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok || path == "" {
		return mcp.NewToolResultError("path is required"), nil
	}

	doc, err := im.GetExampleDocument(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("example not found: %v", err)), nil
	}

	result := map[string]interface{}{
		"id":          doc.ID,
		"category":    doc.Category,
		"subcategory": doc.Subcategory,
		"name":        doc.Name,
		"query":       doc.Query,
	}
	if doc.Variables != "" {
		result["variables"] = doc.Variables
	}
	if doc.Instructions != "" {
		result["instructions"] = doc.Instructions
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to format result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}
