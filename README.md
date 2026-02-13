# twisp-mcp

An MCP server that gives Claude full access to the [Twisp](https://twisp.com) double-entry accounting platform — its GraphQL API, documentation, schema, and 570+ working examples.

Built with the [Model Context Protocol](https://modelcontextprotocol.io), it ships as a single self-contained binary with all content embedded.

## What it does

Twisp MCP lets Claude:

- **Execute GraphQL queries and mutations** against your Twisp instance in real time
- **Search the full GraphQL schema** (406 types, inputs, enums, mutations, queries, scalars, and directives) to discover available API operations
- **Search and read documentation** (94 pages covering accounting concepts, tutorials, guides, and API reference)
- **Search and retrieve working examples** (570+ GraphQL documents demonstrating transaction codes, account setup, multi-currency, ACH, card processing, and more)

This means Claude can look up how Twisp works, find the right mutation, grab a working example, adapt it to your use case, and execute it — all within a single conversation.

## Installation

### Prerequisites

- Go 1.23+

### Build from source

```bash
git clone https://github.com/parsnips/twisp-mcp.git
cd twisp-mcp
make build
```

This produces a `twisp-mcp` binary in the current directory. Optionally install it to `~/bin/`:

```bash
make install
```

### Configure in Claude Code

Add to your Claude Code MCP settings (`~/.claude/settings.json` or project `.claude/settings.json`):

```json
{
  "mcpServers": {
    "twisp": {
      "command": "/path/to/twisp-mcp",
      "env": {
        "TWISP_ENDPOINT": "http://localhost:8080/financial/v1/graphql",
        "TWISP_ACCOUNT_ID": "your-account-id"
      }
    }
  }
}
```

No arguments or index directories needed. The binary is fully self-contained.

## Environment variables

| Variable             | Description                                     | Default                                      |
|----------------------|-------------------------------------------------|----------------------------------------------|
| `TWISP_ENDPOINT`     | GraphQL API endpoint                            | `http://localhost:8080/financial/v1/graphql` |
| `TWISP_ACCOUNT_ID`   | Account ID sent via `X-Twisp-Account-Id` header | `000000000000`                               |
| `TWISP_BEARER_TOKEN` | Bearer token for `Authorization` header         | (none)                                       |

Authentication uses bearer token if set.

## Tools

### `graphql_execute`

Execute any GraphQL query or mutation against the Twisp API.

| Parameter       | Type   | Required | Description                                  |
|-----------------|--------|----------|----------------------------------------------|
| `query`         | string | yes      | GraphQL query or mutation                    |
| `variables`     | object | no       | Variables for the operation                  |
| `operationName` | string | no       | Operation name for multi-operation documents |

### `search_schema`

Full-text search across the Twisp GraphQL schema. Returns type definitions, field lists, and descriptions.

| Parameter | Type    | Required | Description                                                                                                       |
|-----------|---------|----------|-------------------------------------------------------------------------------------------------------------------|
| `query`   | string  | yes      | Search terms (e.g., `"create account"`, `"balance"`)                                                              |
| `kind`    | string  | no       | Filter by element kind: `type`, `input`, `enum`, `mutation`, `query`, `scalar`, `directive`, `interface`, `union` |
| `limit`   | integer | no       | Max results (default 10)                                                                                          |

### `search_docs`

Full-text search across Twisp documentation.

| Parameter  | Type    | Required | Description                                                                                                               |
|------------|---------|----------|---------------------------------------------------------------------------------------------------------------------------|
| `query`    | string  | yes      | Search terms (e.g., `"ACH returns"`, `"velocity controls"`)                                                               |
| `category` | string  | no       | Filter by category: `introduction`, `accounting-core`, `tutorials`, `guides`, `reference`, `infrastructure`, `processors` |
| `limit`    | integer | no       | Max results (default 10)                                                                                                  |

### `get_doc`

Retrieve a full documentation page by path.

| Parameter | Type   | Required | Description                                                                |
|-----------|--------|----------|----------------------------------------------------------------------------|
| `path`    | string | yes      | Doc path, e.g., `"accounting-core/balances"`, `"guides/velocity-controls"` |

### `search_examples`

Full-text search across 570+ GraphQL example documents.

| Parameter  | Type    | Required | Description                                                                |
|------------|---------|----------|----------------------------------------------------------------------------|
| `query`    | string  | yes      | Search terms (e.g., `"post transaction"`, `"multi-currency"`)              |
| `category` | string  | no       | Filter by category: `examples`, `reference`, `fixtures`, `tranCodeLibrary` |
| `limit`    | integer | no       | Max results (default 5)                                                    |

### `get_example`

Retrieve a complete example including GraphQL query, variables, and instructions.

| Parameter | Type   | Required | Description                                                               |
|-----------|--------|----------|---------------------------------------------------------------------------|
| `path`    | string | yes      | Example path, e.g., `"examples/basicTranCodeFlow/001_BasicTranCodeSetup"` |

## How it works

All content (schema files, documentation, and example documents) is embedded into the binary at compile time using Go's `embed` package. At startup, the server parses the GraphQL schema and builds three in-memory [Bleve](https://blevesearch.com) full-text search indexes — typically completing in under 2 seconds. The server then communicates over stdio using the MCP protocol.

## Development

```bash
# Build
make build

# Run tests
make test

# Sync content from twisp/core (requires access to private repo)
TWISP_CORE=/path/to/core make sync
```
