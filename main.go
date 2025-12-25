package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// Parse flags
	kind := flag.String("kind", "", "plugin kind")
	name := flag.String("name", "", "plugin name (repo@version)")
	disableClone := flag.Bool("disable-clone", false, "disable clone")
	sources := flag.String("sources", "", "binary sources")
	flag.Parse()

	_ = kind
	_ = disableClone
	_ = sources

	// Find script path
	scriptPath := findScript(*name)
	if scriptPath == "" {
		fmt.Fprintln(os.Stderr, "Error: posix/script not found")
		fmt.Fprintf(os.Stderr, "Searched with name: %s\n", *name)
		fmt.Fprintf(os.Stderr, "CWD: %s\n", getCwd())
		os.Exit(1)
	}

	fmt.Printf("Executing script: %s\n", scriptPath)

	// Execute the script
	cmd := exec.Command("/bin/sh", scriptPath)
	cmd.Dir = filepath.Dir(filepath.Dir(scriptPath)) // Set working dir to repo root
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

func getCwd() string {
	cwd, _ := os.Getwd()
	return cwd
}

func findScript(name string) string {
	cwd := getCwd()

	// Extract repo path from name (e.g., "github.com/user/repo@refs/tags/v1")
	repoPath := ""
	if name != "" && strings.Contains(name, "@") {
		repoPath = strings.Split(name, "@")[0]
	}

	// Common clone locations
	candidates := []string{
		// Current directory
		filepath.Join(cwd, "posix", "script"),
		// Relative to repo name
		filepath.Join(cwd, repoPath, "posix", "script"),
		// Home directory clone location
		filepath.Join(os.Getenv("HOME"), repoPath, "posix", "script"),
		// Runner clone locations
		filepath.Join("/tmp", repoPath, "posix", "script"),
		filepath.Join("/harness", repoPath, "posix", "script"),
	}

	// Also check DRONE_WORKSPACE if set
	if workspace := os.Getenv("DRONE_WORKSPACE"); workspace != "" {
		candidates = append(candidates, filepath.Join(workspace, "posix", "script"))
		candidates = append(candidates, filepath.Join(workspace, repoPath, "posix", "script"))
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Debug: list current directory
	fmt.Fprintf(os.Stderr, "Debug: Listing CWD contents:\n")
	entries, _ := os.ReadDir(cwd)
	for _, e := range entries {
		fmt.Fprintf(os.Stderr, "  %s\n", e.Name())
	}

	return ""
}
