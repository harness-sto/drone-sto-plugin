package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Find the script relative to the binary location or in current working directory
	scriptPath := findScript()
	if scriptPath == "" {
		fmt.Fprintln(os.Stderr, "Error: posix/script not found")
		os.Exit(1)
	}

	// Execute the script
	cmd := exec.Command("/bin/sh", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error executing script: %v\n", err)
		os.Exit(1)
	}
}

func findScript() string {
	// Try current working directory first
	cwd, _ := os.Getwd()
	candidates := []string{
		filepath.Join(cwd, "posix", "script"),
		filepath.Join(cwd, "posix/script"),
	}

	// Try relative to binary location
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(exeDir, "posix", "script"),
			filepath.Join(exeDir, "..", "posix", "script"),
		)
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}
