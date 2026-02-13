# Twisp MCP Server Makefile
TWISP_CORE ?= $(HOME)/projects/twisp/core

# Source directories in twisp/core
SCHEMA_SRC := $(TWISP_CORE)/services/financial/graphql/schema
DOCS_SRC := $(TWISP_CORE)/services/web/docs/src/pages
EXAMPLES_SRC := $(TWISP_CORE)/services/financial/graphql/documents

# Local directories
EMBED_DIR := embed

.PHONY: all build install test clean sync help

all: build

help:
	@echo "Twisp MCP Server Build System"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build   - Build the twisp-mcp binary (default)"
	@echo "  install - Install to ~/bin/"
	@echo "  test    - Run tests"
	@echo "  clean   - Remove build artifacts"
	@echo "  sync    - Sync content from twisp/core to embed/"
	@echo ""
	@echo "Environment:"
	@echo "  TWISP_CORE - Path to twisp/core repo (default: ~/projects/twisp/core)"

build:
	go build -o twisp-mcp .

install: build
	@mkdir -p ~/bin
	cp twisp-mcp ~/bin/
	@echo "Installed to ~/bin/twisp-mcp"

test:
	go test ./...

clean:
	rm -f twisp-mcp

# Sync embedded content from twisp/core (dev-only)
sync:
	@echo "Syncing schema files..."
	@mkdir -p $(EMBED_DIR)/schema
	@rm -rf $(EMBED_DIR)/schema/*
	cp $(SCHEMA_SRC)/*.graphql $(EMBED_DIR)/schema/
	@echo "Syncing documentation files..."
	@rm -rf $(EMBED_DIR)/docs
	@mkdir -p $(EMBED_DIR)/docs
	cp -r $(DOCS_SRC)/* $(EMBED_DIR)/docs/
	@echo "Syncing example files..."
	@rm -rf $(EMBED_DIR)/examples
	@mkdir -p $(EMBED_DIR)/examples
	@for dir in examples reference fixtures tranCodeLibrary; do \
		if [ -d "$(EXAMPLES_SRC)/$$dir" ]; then \
			cp -r $(EXAMPLES_SRC)/$$dir $(EMBED_DIR)/examples/; \
		fi; \
	done
	@echo "Content synced to $(EMBED_DIR)/"
