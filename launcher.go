package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	binary, err = filepath.Abs(binary)
	if err != nil {
		return err
	}
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find the twisp-mcp executable")
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	// Bazel's generated aws/env launcher changes into its runfiles tree.
	// Restore the user's directory for Codex, keeping the bridge's original
	// cwd and converting credential-file paths before changing directories.
	launchDir, environment := launchEnvironment(cwd, os.Environ())
	args := codexArgs(o, "twisp_bridge_"+strings.ToLower(rand.Text()[:10]), executable, cwd, environment)
	if err := os.Chdir(launchDir); err != nil {
		return fmt.Errorf("cannot return to the original working directory: %w", err)
	}
	// Codex merges configuration tables. A fresh name prevents stale settings
	// (including credentials or a conflicting HTTP transport) from being inherited.
	return replaceProcess(binary, append([]string{binary}, args...), environment)
}

func launchEnvironment(cwd string, environment []string) (string, []string) {
	directory := ""
	for _, entry := range environment {
		if name, value, _ := strings.Cut(entry, "="); name == "BUILD_WORKING_DIRECTORY" {
			directory = value
		}
	}
	if directory == "" || !strings.Contains(filepath.ToSlash(cwd), ".runfiles/") {
		return cwd, environment
	}
	if !filepath.IsAbs(directory) {
		directory = filepath.Join(cwd, directory)
	}
	result := make([]string, 0, len(environment)+1)
	for _, entry := range environment {
		name, value, _ := strings.Cut(entry, "=")
		switch name {
		case "PWD":
			continue
		case "AWS_CONFIG_FILE", "AWS_SHARED_CREDENTIALS_FILE", "AWS_WEB_IDENTITY_TOKEN_FILE", "AWS_CONTAINER_AUTHORIZATION_TOKEN_FILE", "AWS_CA_BUNDLE", "TWISP_MCP_TOKEN_FILE", "SSL_CERT_FILE":
			if value != "" && !filepath.IsAbs(value) {
				value = filepath.Join(cwd, value)
			}
		case "PATH", "SSL_CERT_DIR":
			paths := filepath.SplitList(value)
			for i, path := range paths {
				if path != "" && !filepath.IsAbs(path) {
					paths[i] = filepath.Join(cwd, path)
				}
			}
			value = strings.Join(paths, string(os.PathListSeparator))
		}
		result = append(result, name+"="+value)
	}
	return directory, append(result, "PWD="+directory)
}
