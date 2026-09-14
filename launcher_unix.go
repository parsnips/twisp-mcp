//go:build !windows

package main

import "syscall"

// Replacing the launcher preserves the terminal, signals, and Codex exit code.
func replaceProcess(binary string, args, env []string) error {
	return syscall.Exec(binary, args, env)
}
