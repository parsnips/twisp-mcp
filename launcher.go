package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// JSON strings are also TOML basic strings, including quoting paths and control
// characters. Never interpolate these values into a shell command.
func tomlString(value string) string { data, _ := json.Marshal(value); return string(data) }
func tomlStrings(values []string) string {
	quoted := make([]string, len(values))
	for i, v := range values {
		quoted[i] = tomlString(v)
	}
	return "[" + strings.Join(quoted, ",") + "]"
}

// Forward credential variable NAMES, never their values in argv or config files.
// Includes refreshable AWS container/SSO/profile providers and local GraphQL's
// separate configuration. cwd preserves relative token/config/credential paths.
func codexArgs(o options, name, executable, cwd string, environ []string) []string {
	names := []string{}
	for _, entry := range environ {
		name, _, _ := strings.Cut(entry, "=")
		if strings.HasPrefix(name, "AWS_") || strings.HasPrefix(name, "TWISP_") || name == "HOME" || name == "PATH" || name == "SSL_CERT_FILE" || name == "SSL_CERT_DIR" || strings.HasSuffix(strings.ToUpper(name), "_PROXY") {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	bridgeArgs := []string{"--mode", o.mode, "--account", o.cloud.AccountID, "--env", o.cloud.Environment, "--region", o.cloud.Region, "--url", o.cloud.URL, "serve"}
	config := "mcp_servers." + name + "={command=" + tomlString(executable) + ",args=" + tomlStrings(bridgeArgs) + ",cwd=" + tomlString(cwd) + ",env_vars=" + tomlStrings(names) + ",startup_timeout_sec=45,tool_timeout_sec=130,enabled=true,required=true}"
	args := []string{"-c", config}
	return append(args, o.args...)
}

func launchCodex(o options) error {
	binary, err := exec.LookPath("codex")
	if err != nil {
		return fmt.Errorf("codex is not on PATH; install the Codex CLI first")
	}
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find the twisp-mcp executable")
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	// Codex merges configuration tables. A fresh name prevents stale settings
	// (including credentials or a conflicting HTTP transport) from being inherited.
	name := "twisp_bridge_" + strings.ToLower(rand.Text()[:10])
	return replaceProcess(binary, append([]string{binary}, codexArgs(o, name, executable, cwd, os.Environ())...), os.Environ())
}
