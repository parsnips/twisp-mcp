// Package cloud connects the local MCP to the authenticated Twisp cloud MCP.
package cloud

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
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
	req.Params.ClientInfo = mcp.Implementation{Name: "twisp-mcp-bridge", Version: "3.1.0"}
	if _, err = upstream.Initialize(ctx, req); err != nil {
		_ = upstream.Close()
		return fmt.Errorf("cloud MCP initialization failed; check connectivity and TWISP_MCP credentials")
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
		return nil, fmt.Errorf("cloud MCP tool discovery failed; check connectivity and TWISP_MCP credentials")
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
		return nil, fmt.Errorf("cloud MCP tool call failed; check connectivity and TWISP_MCP credentials")
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
	return http.DefaultTransport.RoundTrip(r)
}

// CheckCredentials is used before launching the agent or a pure cloud server.
func (c *Client) CheckCredentials(ctx context.Context) error {
	if c.config.AccountID == "" {
		return fmt.Errorf("set --account or TWISP_MCP_ACCOUNT_ID")
	}
	_, err := c.tokens.Token(ctx)
	return err
}
