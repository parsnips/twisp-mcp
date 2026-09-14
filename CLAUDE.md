# Twisp MCP development

This is a stdio MCP bridge with two explicit routing modes. In `local` mode
(the default), `graphql_execute` runs directly against `TWISP_ENDPOINT`; other
tools come from the cloud. In `cloud` mode, all tools including GraphQL are
forwarded to hosted MCP. See README.md for configuration and launch commands.

Keep direct GraphQL (`TWISP_*`) and cloud (`TWISP_MCP_*`) credentials separate.
In local mode, never forward GraphQL to cloud MCP or use cloud credentials for
direct GraphQL. Keep local stdio GraphQL usable when cloud discovery fails.
Pure cloud startup and the Codex launcher must preflight cloud access.

Cloud credentials come from a configured token file, explicit token, or AWS
identity exchange, in that order. Never silently change identities on credential
failure. Cache and refresh AWS-derived JWTs in memory before expiry; do not log
tokens, put them in argv/config files, or automatically replay failed tool calls.
AWS exchange must only target the selected Twisp deployment.

Go 1.25.5+ is required. MCP SDK: github.com/mark3labs/mcp-go v1.0.0.
There is no embedded corpus, Bleve index, Bedrock client, or content sync step.
Builds and tests must not require AWS credentials or a live cloud service.

Validation:

```bash
make build
go test -race ./...
go vet ./...
```

Use `GOWORK=off` if an inherited Go workspace excludes this repository.

Cloud tools are discovered rather than hardcoded. Preserve cloud tool schemas,
annotations, results, and metadata. Refresh registration only when the catalog
changes, to avoid a tools/list notification loop. Keep transport headers and
endpoints under local configuration control, never tool arguments, and do not
follow HTTP redirects. Preserve terminal behavior and child exit status in the
launcher. Its session configuration must not modify global Codex settings.
