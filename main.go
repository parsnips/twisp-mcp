package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/parsnips/twisp-mcp/internal/cloud"
	"github.com/parsnips/twisp-mcp/internal/graphql"
	"github.com/parsnips/twisp-mcp/internal/server"
)

type options struct {
	cloud   cloud.Config
	mode    string
	command string
	args    []string
	version bool
}

func parseOptions(args []string, stderr io.Writer) (options, error) {
	o := options{cloud: cloud.ConfigFromEnvironment(), mode: os.Getenv("TWISP_MCP_MODE"), command: "serve"}
	if o.mode == "" {
		o.mode = "local"
	}
	flags := flag.NewFlagSet("twisp-mcp", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.BoolVar(&o.version, "version", false, "print version and exit")
	flags.StringVar(&o.mode, "mode", o.mode, "tool routing: local or cloud")
	flags.StringVar(&o.cloud.AccountID, "account", o.cloud.AccountID, "cloud Twisp tenant/account (or TWISP_MCP_ACCOUNT_ID)")
	flags.StringVar(&o.cloud.Environment, "env", o.cloud.Environment, "Twisp deployment environment (or TWISP_MCP_ENV / ZONE)")
	flags.StringVar(&o.cloud.Region, "region", o.cloud.Region, "Twisp deployment region (or TWISP_MCP_REGION / AWS_REGION)")
	endpoint := os.Getenv("TWISP_MCP_URL")
	flags.StringVar(&endpoint, "url", endpoint, "override cloud MCP URL; custom URLs require explicit tokens")
	flags.Usage = func() {
		fmt.Fprintln(stderr, "Usage: twisp-mcp [options] [serve | codex [codex arguments...]]")
		fmt.Fprintln(stderr, "Options must precede the command. No command means serve MCP over stdio.")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		return o, err
	}
	if o.version {
		return o, nil
	}
	if o.mode != "local" && o.mode != "cloud" {
		return o, fmt.Errorf("--mode must be local or cloud")
	}
	o.cloud.ForwardGraphQL = o.mode == "cloud"
	if endpoint == "" {
		endpoint = cloud.Endpoint(o.cloud.Environment, o.cloud.Region)
	}
	o.cloud.URL = endpoint
	if remaining := flags.Args(); len(remaining) > 0 {
		o.command, o.args = remaining[0], remaining[1:]
	}
	if o.command != "serve" && o.command != "codex" {
		return o, fmt.Errorf("unknown command %q; expected serve or codex", o.command)
	}
	if o.command == "serve" && len(o.args) > 0 {
		return o, fmt.Errorf("serve accepts no positional arguments; place options before serve")
	}
	return o, nil
}

func run(args []string) error {
	o, err := parseOptions(args, os.Stderr)
	if err != nil {
		return err
	}
	if o.version {
		fmt.Println(server.Version)
		return nil
	}
	upstream, err := cloud.New(o.cloud)
	if err != nil {
		return err
	}
	defer upstream.Close()
	if o.command == "codex" {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := upstream.CheckCredentials(ctx); err != nil {
			return err
		}
		if _, err := upstream.ListTools(ctx); err != nil {
			return fmt.Errorf("cloud preflight: %w", err)
		}
		upstream.Close()
		cancel()
		return launchCodex(o)
	}
	var local *graphql.Client
	if o.mode == "local" {
		local = graphql.NewClient()
	}
	s := server.New(local, upstream)
	if o.mode == "cloud" {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := upstream.CheckCredentials(ctx); err != nil {
			return err
		}
		if err := s.CheckCloud(ctx); err != nil {
			return fmt.Errorf("cloud startup: %w", err)
		}
		cancel()
	}
	return s.Serve()
}

func main() {
	if err := run(os.Args[1:]); err != nil && !errors.Is(err, flag.ErrHelp) {
		log.Fatal(err)
	}
}
