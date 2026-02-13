# Twisp MCP Server

MCP (Model Context Protocol) server for interacting with the Twisp double-entry accounting API.

## Features

- **Generic GraphQL execution** - Execute any GraphQL query or mutation against the Twisp API
- **Full-text search** - Search schema, documentation, and examples using Bleve (in-memory indexes)
- **Self-contained binary** - Schema, docs, and 585+ examples embedded via `go:embed`, indexes built at startup

## Building

```bash
# Build the binary
make build

# Install to ~/bin/
make install

# Clean build artifacts
make clean
```

No external dependencies or index directories needed — the binary is fully self-contained.

### Prerequisites

- Go 1.23+

### Syncing Content

To update embedded content from `twisp/core` (requires access to private repo):

```bash
# Default: ~/projects/twisp/core
make sync

# Or specify a custom path
TWISP_CORE=/path/to/core make sync
```

After syncing, rebuild with `make build`.

## Usage

### Running the server

```bash
./twisp-mcp
```

Indexes are built in-memory at startup (~1-3 seconds).

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TWISP_ENDPOINT` | GraphQL API endpoint | `http://localhost:8080/graphql` |
| `TWISP_ACCOUNT_ID` | Account ID header | `000000000000` |
| `TWISP_API_KEY` | API key for authentication | (none) |
| `TWISP_BEARER_TOKEN` | Bearer token for authentication | (none) |

### Claude Code Configuration

Add to `~/.config/claude-code/config.json`:

```json
{
  "mcpServers": {
    "twisp": {
      "command": "/path/to/twisp-mcp",
      "env": {
        "TWISP_ENDPOINT": "http://localhost:8080/graphql",
        "TWISP_ACCOUNT_ID": "your-account-id"
      }
    }
  }
}
```

## MCP Tools

### graphql_execute

Execute arbitrary GraphQL queries or mutations against the Twisp API.

**Parameters:**
- `query` (required): GraphQL query or mutation string
- `variables` (optional): JSON object of variables
- `operationName` (optional): Operation name for multi-operation documents

### search_schema

Search the GraphQL schema for types, mutations, queries, and fields.

**Parameters:**
- `query` (required): Search terms
- `kind` (optional): Filter by type/input/enum/mutation/query/scalar/directive
- `limit` (optional): Max results (default 10)

### search_docs

Search Twisp documentation.

**Parameters:**
- `query` (required): Search terms
- `category` (optional): introduction/accounting-core/tutorials/guides/reference/infrastructure
- `limit` (optional): Max results (default 10)

### get_doc

Retrieve full documentation page by path.

**Parameters:**
- `path` (required): Doc path (e.g., "accounting-core/balances")

### search_examples

Search 585 GraphQL example documents.

**Parameters:**
- `query` (required): Search terms
- `category` (optional): examples/reference/fixtures/tranCodeLibrary
- `limit` (optional): Max results (default 5)

### get_example

Retrieve complete example by path.

**Parameters:**
- `path` (required): Example path (e.g., "examples/basicTranCodeFlow/001_BasicTranCodeSetup")

## Project Structure

```
twisp-mcp/
├── main.go                     # Entry point: embeds content, builds indexes, starts server
├── embed.go                    # go:embed directive for content
├── go.mod
├── Makefile                    # build, install, sync, test
├── CLAUDE.md                   # This file
│
├── internal/
│   ├── server/server.go        # MCP server setup and tool registration
│   ├── graphql/client.go       # HTTP client with Execute method
│   ├── search/
│   │   ├── types.go            # Document types for indexing
│   │   ├── index.go            # IndexManager, search/get methods, Bleve mappings
│   │   └── loader.go           # BuildIndexesFromFS: parses content, builds in-memory indexes
│   └── tools/
│       ├── graphql_execute.go  # Generic GraphQL executor
│       ├── search_schema.go    # Schema search tool
│       ├── search_docs.go      # Documentation search tool
│       └── search_examples.go  # Examples search tool
│
└── embed/                      # Embedded content (checked into git)
    ├── schema/                 # GraphQL schema files (.graphql)
    ├── docs/                   # Markdown documentation
    └── examples/               # Example GraphQL documents
```

## Development

```bash
# Run tests
make test

# Sync content from twisp/core and rebuild
make sync && make build
```
