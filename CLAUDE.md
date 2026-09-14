# Twisp MCP development

This is a local stdio MCP tool proxy. `graphql_execute` runs directly against
`TWISP_ENDPOINT`; all other tool definitions and calls come from the cloud MCP.
See README.md for configuration and the interaction diagram.

Keep GraphQL (`TWISP_*`) and cloud MCP (`TWISP_MCP_*`) credentials separate.
Never forward GraphQL calls to cloud MCP, use cloud credentials for GraphQL,
accept endpoint overrides from tool arguments, or log tokens. Keep local
GraphQL usable when cloud discovery or cloud tool calls fail.

Go 1.25.5+ is required. MCP SDK: github.com/mark3labs/mcp-go v1.0.0.
There is no embedded corpus, Bleve index, Bedrock client, or content sync step.
Builds and tests must not require AWS credentials or a live cloud service.

Validation:

```bash
make build
go test -race ./...
go vet ./...
```

Cloud tools are discovered rather than hardcoded. Preserve cloud tool schemas,
annotations, results, and metadata. Refresh registration only when the catalog
changes, to avoid a tools/list notification loop. Keep cloud transport headers
under local configuration control and do not follow HTTP redirects.
