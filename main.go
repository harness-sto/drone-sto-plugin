package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
)

//go:embed sto-plugin-binary
var stoBinary []byte

func main() {
	tmp := "/tmp/sto-plugin-embedded"
	if err := os.WriteFile(tmp, stoBinary, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command(tmp, "--run-strategy", "single-container")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}
