.PHONY: all build install test clean help

all: build

build:
	go build -trimpath -o twisp-mcp .

install: build
	mkdir -p $(HOME)/bin
	cp twisp-mcp $(HOME)/bin/

test:
	go test -race ./...
	go vet ./...

clean:
	rm -f twisp-mcp

help:
	@echo "make build    Build the local MCP proxy"
	@echo "make install  Install into ~/bin"
	@echo "make test     Run tests with the race detector and go vet"
	@echo "make clean    Remove the binary"
