// Package cloud connects the local MCP to the authenticated Twisp cloud MCP.
package cloud

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

const DefaultURL = "https://api.us-east-1.dev.twisp.com/mcp"

type Config struct {
	URL, AccountID, BearerToken, TokenFile string
	Environment, Region                    string
	ForwardGraphQL                         bool
}

func ConfigFromEnvironment() Config {
	environment := first(os.Getenv("TWISP_MCP_ENV"), os.Getenv("ZONE"), "dev")
	region := first(os.Getenv("TWISP_MCP_REGION"), os.Getenv("AWS_REGION"), os.Getenv("AWS_DEFAULT_REGION"), "us-east-1")
	endpoint := os.Getenv("TWISP_MCP_URL")
	if endpoint == "" {
		endpoint = Endpoint(environment, region)
	}
	return Config{URL: endpoint, AccountID: os.Getenv("TWISP_MCP_ACCOUNT_ID"), BearerToken: os.Getenv("TWISP_MCP_BEARER_TOKEN"), TokenFile: os.Getenv("TWISP_MCP_TOKEN_FILE"), Environment: environment, Region: region}
}

func first(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func Endpoint(environment, region string) string {
	return fmt.Sprintf("https://api.%s.%s.twisp.com/mcp", region, environment)
}

// Client initializes lazily, allowing local GraphQL to work without cloud access.
type Client struct {
	config Config
	gate   chan struct{}
	client *client.Client
	tokens *tokenSource
}

func New(config Config) (*Client, error) {
	u, err := url.Parse(config.URL)
	if err != nil || u.Hostname() == "" || u.User != nil || u.Fragment != "" {
		return nil, fmt.Errorf("invalid TWISP_MCP_URL")
	}
	loopback := u.Hostname() == "localhost" || net.ParseIP(u.Hostname()).IsLoopback()
	if u.Scheme != "https" && !(u.Scheme == "http" && loopback) {
		return nil, fmt.Errorf("TWISP_MCP_URL requires HTTPS (HTTP is allowed on loopback)")
	}
	return &Client{config: config, tokens: &tokenSource{config: config, now: time.Now}, gate: make(chan struct{}, 1)}, nil
}

func (c *Client) connect(ctx context.Context) error {
	if c.client != nil {
		return nil
	}
	httpClient := &http.Client{Timeout: 2 * time.Minute, Transport: credentialTransport{config: c.config, tokens: c.tokens}, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	upstream, err := client.NewStreamableHttpClient(c.config.URL, transport.WithHTTPBasicClient(httpClient))
	if err != nil {
		return fmt.Errorf("cannot create cloud MCP connection")
	}
	if err = upstream.Start(ctx); err != nil {
		_ = upstream.Close()
		return fmt.Errorf("cannot start cloud MCP connection")
	}
	req := mcp.InitializeRequest{}
	// Match the hosted MCP's Streamable HTTP handshake rather than probing a
	// newer discovery protocol on the API proxy.
	req.Params.ProtocolVersion = "2025-03-26"
	req.Params.ClientInfo = mcp.Implementation{Name: "twisp-mcp-bridge", Version: "3.1.1"}
	if _, err = upstream.Initialize(ctx, req); err != nil {
		_ = upstream.Close()
		return connectionError("initialization", err)
	}
	c.client = upstream
	return nil
}

func (c *Client) ListTools(ctx context.Context) (*mcp.ListToolsResult, error) {
	if err := c.acquire(ctx); err != nil {
		return nil, err
	}
	defer c.release()
	if err := c.connect(ctx); err != nil {
		return nil, err
	}
	result, err := c.client.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		c.reset()
		return nil, connectionError("tool discovery", err)
	}
	return result, nil
}

func (c *Client) CallTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Local mode keeps GraphQL entirely separate from cloud credentials.
	if request.Params.Name == "graphql_execute" && !c.config.ForwardGraphQL {
		return nil, fmt.Errorf("graphql_execute is local only")
	}
	if err := c.acquire(ctx); err != nil {
		return nil, err
	}
	defer c.release()
	if err := c.connect(ctx); err != nil {
		return nil, err
	}
	// Only forward the tool parameters. Inbound transport headers must not
	// override the separately configured cloud identity.
	outbound := mcp.CallToolRequest{}
	outbound.Params = request.Params
	result, err := c.client.CallTool(ctx, outbound)
	if err != nil {
		c.reset()
		return nil, connectionError("tool call", err)
	}
	return result, nil
}

func (c *Client) reset() {
	if c.client != nil {
		_ = c.client.Close()
		c.client = nil
	}
}
func (c *Client) Close() { _ = c.acquire(context.Background()); defer c.release(); c.reset() }
func (c *Client) acquire(ctx context.Context) error {
	select {
	case c.gate <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (c *Client) release() { <-c.gate }

type credentialTransport struct {
	config Config
	tokens *tokenSource
}

func (t credentialTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.config.AccountID == "" {
		return nil, fmt.Errorf("set --account or TWISP_MCP_ACCOUNT_ID")
	}
	token, err := t.tokens.Token(r.Context())
	if err != nil {
		return nil, err
	}
	r = r.Clone(r.Context())
	r.Header.Set("Authorization", "Bearer "+token)
	r.Header.Set("X-Twisp-Account-Id", t.config.AccountID)
	response, err := http.DefaultTransport.RoundTrip(r)
	if err == nil && response.StatusCode >= 300 {
		// Recognize the ALB's fixed rejection without exposing arbitrary response
		// bodies (which could contain credentials). Keep diagnostics bounded.
		body, _ := io.ReadAll(io.LimitReader(response.Body, 256))
		_ = response.Body.Close()
		return nil, &httpStatusError{status: response.StatusCode, defaultRejection: response.StatusCode == 401 && strings.TrimSpace(string(body)) == "unauthorized"}
	}
	return response, err
}

// Intercept HTTP failures before the SDK maps initialize's 4xx responses to
// "legacy SSE" or includes an untrusted response body in its error text.
type httpStatusError struct {
	status           int
	defaultRejection bool
}

func (e *httpStatusError) Error() string {
	hint := "check the cloud MCP endpoint"
	switch e.status {
	case 400:
		hint = "check --account; tenant aliases require the alias/ prefix"
	case 401:
		hint = "the cloud API rejected the bearer token; check token expiry and issuer"
		if e.defaultRejection {
			hint = "plain unauthorized response; check that /mcp is deployed and routed by the load balancer"
		}
	case 403:
		hint = "access denied; check --account (aliases use alias/) and the AWS identity's Twisp client registration"
	case 404, 405:
		hint = "check that hosted MCP is deployed at the selected URL"
	case 502, 503, 504:
		hint = "the cloud MCP service or runtime is unavailable"
	}
	return fmt.Sprintf("HTTP %d: %s", e.status, hint)
}

func connectionError(operation string, err error) error {
	var status *httpStatusError
	if errors.As(err, &status) {
		return fmt.Errorf("cloud MCP %s failed: %w", operation, status)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("cloud MCP %s timed out", operation)
	}
	if errors.Is(err, context.Canceled) {
		return fmt.Errorf("cloud MCP %s canceled", operation)
	}
	return fmt.Errorf("cloud MCP %s failed; check connectivity, protocol compatibility, and TWISP_MCP credentials", operation)
}

// CheckCredentials is used before launching the agent or a pure cloud server.
func (c *Client) CheckCredentials(ctx context.Context) error {
	if c.config.AccountID == "" {
		return fmt.Errorf("set --account or TWISP_MCP_ACCOUNT_ID")
	}
	_, err := c.tokens.Token(ctx)
	return err
}
