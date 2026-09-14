package main

import (
	"context"
	"github.com/mark3labs/mcp-go/client"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	sdk "github.com/mark3labs/mcp-go/server"
)

func TestLauncherProcess(t *testing.T) {
	if os.Getenv("TWISP_TEST_HELPER") != "1" {
		return
	}
	for i, arg := range os.Args {
		if arg == "--" {
			if err := run(os.Args[i+1:]); err != nil {
				os.Exit(23)
			}
			os.Exit(0)
		}
	}
	os.Exit(24)
}

func TestCodexLauncherPreflightAndExec(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell fixture requires Unix")
	}
	remote := sdk.NewMCPServer("fixture", "1")
	remote.AddTool(mcp.NewTool("graphql_execute"), func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("ok"), nil
	})
	transport := sdk.NewStreamableHTTPServer(remote, sdk.WithStateLess(true))
	host := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer fake-jwt" || r.Header.Get("X-Twisp-Account-Id") != "tenant" {
			t.Error("wrong preflight identity")
		}
		transport.ServeHTTP(w, r)
	}))
	defer host.Close()
	dir := filepath.Join(t.TempDir(), "path with spaces")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	capture := filepath.Join(dir, "args")
	script := "#!/bin/sh\n[ \"$AWS_SESSION_TOKEN\" = fake-secret ] || exit 25\nprintf '%s\\n' \"$@\" > \"$TWISP_TEST_CAPTURE\"\nexit 17\n"
	if err := os.WriteFile(filepath.Join(dir, "codex"), []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	// Explicitly isolate the child's configuration from the developer's shell.
	env := []string{"PATH=" + dir + string(os.PathListSeparator) + os.Getenv("PATH"), "HOME=" + dir, "TWISP_TEST_HELPER=1", "TWISP_TEST_CAPTURE=" + capture, "TWISP_MCP_BEARER_TOKEN=fake-jwt", "AWS_SESSION_TOKEN=fake-secret"}
	command := func(endpoint string) *exec.Cmd {
		cmd := exec.Command(os.Args[0], "-test.run=^TestLauncherProcess$", "--", "--mode", "cloud", "--env", "dev", "--region", "us-east-1", "--account", "tenant", "--url", endpoint, "codex", "exec", "a prompt with spaces")
		cmd.Env = env
		cmd.Dir = dir
		return cmd
	}
	output, err := command(host.URL).CombinedOutput()
	if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 17 {
		t.Fatalf("launcher lost exit code: %v %s", err, output)
	}
	data, err := os.ReadFile(capture)
	if err != nil {
		t.Fatal(err)
	}
	args := string(data)
	for _, expected := range []string{"mcp_servers.twisp_bridge_", `"--mode","cloud"`, `"AWS_SESSION_TOKEN"`, `"TWISP_MCP_BEARER_TOKEN"`, "\nexec\na prompt with spaces\n", "required=true"} {
		if !strings.Contains(args, expected) {
			t.Fatalf("missing %q in launcher arguments", expected)
		}
	}
	if strings.Contains(args, "fake-secret") || strings.Contains(args, "fake-jwt") {
		t.Fatal("secret leaked to argv")
	}
	if _, err := os.Stat(filepath.Join(dir, ".codex", "config.toml")); !os.IsNotExist(err) {
		t.Fatal("launcher wrote global config")
	}

	// Reproduce the Bazel wrapper: execute from runfiles with a relative AWS
	// config path, but launch Codex in BUILD_WORKING_DIRECTORY.
	runfiles := filepath.Join(dir, "vault.runfiles", "_main")
	workspace := filepath.Join(dir, "workspace")
	for _, path := range []string{runfiles, workspace} {
		if err := os.MkdirAll(path, 0700); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(runfiles, "awsconfig"), []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	wrapperScript := "#!/bin/sh\n[ -f \"$AWS_CONFIG_FILE\" ] || exit 26\n[ \"$PWD\" = \"$BUILD_WORKING_DIRECTORY\" ] || exit 27\npwd > \"$TWISP_TEST_CAPTURE\"\nexit 17\n"
	if err := os.WriteFile(filepath.Join(dir, "codex"), []byte(wrapperScript), 0700); err != nil {
		t.Fatal(err)
	}
	wrapped := command(host.URL)
	wrapped.Dir = runfiles
	wrapped.Env = append(wrapped.Env, "BUILD_WORKING_DIRECTORY="+workspace, "AWS_CONFIG_FILE=awsconfig")
	if output, err := wrapped.CombinedOutput(); err == nil {
		t.Fatal("lost wrapped exit code")
	} else if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 17 {
		t.Fatalf("wrapped launch failed: %v %s", err, output)
	}
	launched, err := os.ReadFile(capture)
	if err != nil {
		t.Fatal(err)
	}
	resolvedWorkspace, err := filepath.EvalSymlinks(workspace)
	if err != nil {
		t.Fatal(err)
	}
	actualWorkspace, err := filepath.EvalSymlinks(strings.TrimSpace(string(launched)))
	if err != nil {
		t.Fatal(err)
	}
	if actualWorkspace != resolvedWorkspace {
		t.Fatalf("wrong working directory: %q", launched)
	}
	if err := os.Remove(capture); err != nil {
		t.Fatal(err)
	}
	rejected := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "unauthorized", 401) }))
	defer rejected.Close()
	if err := command(rejected.URL).Run(); err == nil {
		t.Fatal("launched despite failed preflight")
	}
	if _, err := os.Stat(capture); !os.IsNotExist(err) {
		t.Fatal("Codex ran after rejected preflight")
	}
}

func TestOptionsAndWrapperEnvironment(t *testing.T) {
	for _, name := range []string{"TWISP_MCP_MODE", "TWISP_MCP_ENV", "TWISP_MCP_REGION", "TWISP_MCP_URL", "TWISP_MCP_ACCOUNT_ID", "AWS_DEFAULT_REGION"} {
		t.Setenv(name, "")
	}
	t.Setenv("ZONE", "dev")
	t.Setenv("AWS_REGION", "us-east-1")
	o, err := parseOptions([]string{"--mode", "cloud", "--account", "tenant", "codex", "--model", "test"}, io.Discard)
	if err != nil || o.cloud.URL != "https://api.us-east-1.dev.twisp.com/mcp" || !o.cloud.ForwardGraphQL || o.command != "codex" || len(o.args) != 2 {
		t.Fatal(o, err)
	}
	o, err = parseOptions([]string{"--env", "cloud", "--region", "us-west-2"}, io.Discard)
	if err != nil || o.mode != "local" || o.cloud.URL != "https://api.us-west-2.cloud.twisp.com/mcp" {
		t.Fatal(o, err)
	}
	for _, args := range [][]string{{"--mode", "wrong"}, {"serve", "--mode", "cloud"}, {"claude"}} {
		if _, err := parseOptions(args, io.Discard); err == nil {
			t.Fatal("accepted invalid arguments", args)
		}
	}
}

// Exercise the real serve path after its startup context has been canceled.
func TestCloudStdioProcess(t *testing.T) {
	var calls atomic.Int32
	remote := sdk.NewMCPServer("fixture", "1")
	remote.AddTool(mcp.NewTool("graphql_execute", mcp.WithString("query", mcp.Required())), func(_ context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		calls.Add(1)
		if r.GetString("query", "") != "query { __typename }" {
			t.Error("changed query")
		}
		return mcp.NewToolResultText("hosted"), nil
	})
	transport := sdk.NewStreamableHTTPServer(remote, sdk.WithStateLess(true))
	host := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer fixture-token" || r.Header.Get("X-Twisp-Account-Id") != "tenant" {
			t.Error("wrong stdio bridge identity")
		}
		transport.ServeHTTP(w, r)
	}))
	defer host.Close()
	c, err := client.NewStdioMCPClient(os.Args[0], []string{"TWISP_TEST_HELPER=1", "TWISP_MCP_TOKEN_FILE=", "TWISP_MCP_BEARER_TOKEN=fixture-token"}, "-test.run=^TestLauncherProcess$", "--", "--mode", "cloud", "--account", "tenant", "--url", host.URL, "serve")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	init := mcp.InitializeRequest{}
	init.Params.ProtocolVersion = "2025-03-26"
	if _, err := c.Initialize(ctx, init); err != nil {
		t.Fatal(err)
	}
	list, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil || len(list.Tools) != 1 || list.Tools[0].Name != "graphql_execute" {
		t.Fatal(list, err)
	}
	call := mcp.CallToolRequest{}
	call.Params.Name = "graphql_execute"
	call.Params.Arguments = map[string]any{"query": "query { __typename }"}
	result, err := c.CallTool(ctx, call)
	if err != nil || result.IsError || calls.Load() != 1 {
		t.Fatal(result, err)
	}
}
