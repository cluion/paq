package provider

import (
	"context"
	"os/exec"
	"time"
)

// CommandRunner executes a system command and returns its stdout.
type CommandRunner func(name string, args ...string) (string, error)

// defaultRunner executes system commands with a 30-second timeout.
var defaultRunner CommandRunner = func(name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.Output()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}

// commandExists checks whether a command is available in PATH.
func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
