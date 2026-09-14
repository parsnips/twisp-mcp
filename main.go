package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/parsnips/twisp-mcp/internal/cloud"
	"github.com/parsnips/twisp-mcp/internal/graphql"
	"github.com/parsnips/twisp-mcp/internal/server"
)

func main() {
	version := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *version {
		fmt.Println(server.Version)
		return
	}
	upstream, err := cloud.New(cloud.ConfigFromEnvironment())
	if err != nil {
		log.Fatal(err)
	}
	defer upstream.Close()
	if err := server.New(graphql.NewClient(), upstream).Serve(); err != nil {
		log.Fatal(err)
	}
}
