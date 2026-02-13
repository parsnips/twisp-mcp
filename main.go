package main

import (
	"io/fs"
	"log"

	"github.com/parsnips/twisp-mcp/internal/search"
	"github.com/parsnips/twisp-mcp/internal/server"
)

func main() {
	// Strip the "embed" prefix so contentFS has schema/, docs/, examples/ at root
	contentFS, err := fs.Sub(embeddedContent, "embed")
	if err != nil {
		log.Fatalf("Failed to access embedded content: %v", err)
	}

	// Build in-memory indexes from embedded content
	im, err := search.BuildIndexesFromFS(contentFS)
	if err != nil {
		log.Fatalf("Failed to build indexes: %v", err)
	}
	defer im.Close()

	// Log index stats
	if count, err := im.SchemaStats(); err == nil {
		log.Printf("Schema index: %d documents", count)
	}
	if count, err := im.DocsStats(); err == nil {
		log.Printf("Docs index: %d documents", count)
	}
	if count, err := im.ExamplesStats(); err == nil {
		log.Printf("Examples index: %d documents", count)
	}

	// Create and start the server
	srv, err := server.New(im)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Serve(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
