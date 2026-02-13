package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Client is a GraphQL client for the Twisp API.
type Client struct {
	endpoint    string
	accountID   string
	apiKey      string
	bearerToken string
	httpClient  *http.Client
}

// GraphQLRequest represents a GraphQL request payload.
type GraphQLRequest struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	OperationName string                 `json:"operationName,omitempty"`
}

// GraphQLResponse represents a GraphQL response.
type GraphQLResponse struct {
	Data   json.RawMessage `json:"data,omitempty"`
	Errors []GraphQLError  `json:"errors,omitempty"`
}

// GraphQLError represents a GraphQL error.
type GraphQLError struct {
	Message    string                 `json:"message"`
	Path       []interface{}          `json:"path,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// NewClient creates a new GraphQL client from environment variables.
func NewClient() *Client {
	endpoint := os.Getenv("TWISP_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8080/financial/v1/graphql"
	}

	accountID := os.Getenv("TWISP_ACCOUNT_ID")
	if accountID == "" {
		accountID = "000000000000"
	}

	return &Client{
		endpoint:    endpoint,
		accountID:   accountID,
		apiKey:      os.Getenv("TWISP_API_KEY"),
		bearerToken: os.Getenv("TWISP_BEARER_TOKEN"),
		httpClient:  &http.Client{},
	}
}

// NewClientWithConfig creates a new GraphQL client with explicit configuration.
func NewClientWithConfig(endpoint, accountID, apiKey, bearerToken string) *Client {
	if endpoint == "" {
		endpoint = "http://localhost:8080/financial/v1/graphql"
	}
	if accountID == "" {
		accountID = "000000000000"
	}
	return &Client{
		endpoint:    endpoint,
		accountID:   accountID,
		apiKey:      apiKey,
		bearerToken: bearerToken,
		httpClient:  &http.Client{},
	}
}

// Execute executes a GraphQL query or mutation and returns the raw response.
func (c *Client) Execute(ctx context.Context, query string, variables map[string]interface{}, operationName string) (*GraphQLResponse, error) {
	reqBody := GraphQLRequest{
		Query:         query,
		Variables:     variables,
		OperationName: operationName,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Twisp-Account-Id", c.accountID)

	// Add authentication headers
	if c.bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	} else if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(respBody))
	}

	var gqlResp GraphQLResponse
	if err := json.Unmarshal(respBody, &gqlResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &gqlResp, nil
}

// ExecuteRaw executes a GraphQL query and returns the response as a formatted JSON string.
func (c *Client) ExecuteRaw(ctx context.Context, query string, variables map[string]interface{}, operationName string) (string, error) {
	resp, err := c.Execute(ctx, query, variables, operationName)
	if err != nil {
		return "", err
	}

	// Format the response nicely
	result := make(map[string]interface{})

	if len(resp.Errors) > 0 {
		errors := make([]map[string]interface{}, len(resp.Errors))
		for i, e := range resp.Errors {
			errors[i] = map[string]interface{}{
				"message": e.Message,
			}
			if e.Path != nil {
				errors[i]["path"] = e.Path
			}
			if e.Extensions != nil {
				errors[i]["extensions"] = e.Extensions
			}
		}
		result["errors"] = errors
	}

	if resp.Data != nil {
		var data interface{}
		if err := json.Unmarshal(resp.Data, &data); err == nil {
			result["data"] = data
		}
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(output), nil
}
