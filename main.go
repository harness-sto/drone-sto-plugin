package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
)

//go:embed posix/script
var script string

func main() {
	// Write embedded script to temp file
	tmpFile := "/tmp/sto-script.sh"
	if err := os.WriteFile(tmpFile, []byte(script), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing script: %v\n", err)
		os.Exit(1)
	}

	// Execute the script
	cmd := exec.Command("/bin/sh", tmpFile)
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
