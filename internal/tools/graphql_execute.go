package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/parsnips/twisp-mcp/internal/graphql"
)

// GraphQLExecuteTool returns the tool definition and handler for graphql_execute.
func GraphQLExecuteTool() (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.Tool{
		Name:        "graphql_execute",
		Description: "Execute arbitrary GraphQL queries or mutations against the Twisp API. Authentication is configured via environment variables: TWISP_ENDPOINT (default: http://localhost:8080/financial/v1/graphql), TWISP_ACCOUNT_ID (default: 000000000000), and optionally TWISP_API_KEY or TWISP_BEARER_TOKEN.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "The GraphQL query or mutation string to execute",
				},
				"variables": map[string]interface{}{
					"type":        "object",
					"description": "Optional JSON object of variables for the query/mutation",
				},
				"operationName": map[string]interface{}{
					"type":        "string",
					"description": "Optional operation name for multi-operation documents",
				},
			},
			Required: []string{"query"},
		},
	}

	return tool, handleGraphQLExecute
}

func handleGraphQLExecute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract query
	query, ok := request.Params.Arguments["query"].(string)
	if !ok || query == "" {
		return mcp.NewToolResultError("query is required and must be a string"), nil
	}

	// Extract optional variables
	var variables map[string]interface{}
	if vars, ok := request.Params.Arguments["variables"]; ok && vars != nil {
		switch v := vars.(type) {
		case map[string]interface{}:
			variables = v
		case string:
			// Try to parse as JSON string
			if err := json.Unmarshal([]byte(v), &variables); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to parse variables as JSON: %v", err)), nil
			}
		}
	}

	// Extract optional operation name
	var operationName string
	if opName, ok := request.Params.Arguments["operationName"].(string); ok {
		operationName = opName
	}

	// Execute the query
	client := graphql.NewClient()
	result, err := client.ExecuteRaw(ctx, query, variables, operationName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("GraphQL execution failed: %v", err)), nil
	}

	return mcp.NewToolResultText(result), nil
}
