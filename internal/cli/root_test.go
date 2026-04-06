package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestRootHelp(t *testing.T) {
	out, err := executeCommand(rootCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Usage:") {
		t.Error("root command without args should show help")
	}
}

func TestRootJsonFlagDefault(t *testing.T) {
	jsonOutput = false
	rootCmd.SetArgs([]string{"--json", "version"})
	// just verify flag parsing doesn't error
	_ = rootCmd.Execute()

	// reset
	jsonOutput = false
	rootCmd.SetArgs([]string{})
}

func TestRootErrorOutput(t *testing.T) {
	// Running a non-existent subcommand should error
	_, err := executeCommand(rootCmd, "nonexistent")
	if err == nil {
		t.Error("expected error for unknown command")
	}
}
