package utils

import (
	"context"
	"os/exec"
	"time"
)

// ExecuteWithSudo runs a command with sudo privileges
func ExecuteWithSudo(timeout time.Duration, command string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Prepend sudo to the command
	cmdArgs := append([]string{command}, args...)
	cmd := exec.CommandContext(ctx, "sudo", cmdArgs...)

	// Execute and capture output
	output, err := cmd.CombinedOutput()
	return output, err
}
