package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cluion/paq/internal/provider"
)

type listTestProvider struct {
	name        string
	displayName string
	detect      bool
}

func (p *listTestProvider) Name() string                    { return p.name }
func (p *listTestProvider) DisplayName() string              { return p.displayName }
func (p *listTestProvider) Detect() bool                     { return p.detect }
func (p *listTestProvider) List() ([]provider.Package, error) { return nil, nil }

func TestListTable(t *testing.T) {
	provider.Register(&listTestProvider{name: "testlist-brew", displayName: "Homebrew", detect: true})
	provider.Register(&listTestProvider{name: "testlist-snap", displayName: "Snap", detect: false})

	RegisterProviderCommands()

	jsonOutput = false
	out, err := executeCommand(rootCmd, "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "Homebrew") {
		t.Error("output should contain 'Homebrew'")
	}
	if !strings.Contains(out, "available") {
		t.Error("output should contain 'available'")
	}
	if !strings.Contains(out, "not found") {
		t.Error("output should contain 'not found'")
	}
}

func TestListJSON(t *testing.T) {
	jsonOutput = true
	rootCmd.SetArgs([]string{"list"})
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	err := rootCmd.Execute()
	jsonOutput = false

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"sources"`) {
		t.Error("JSON output should contain 'sources' key")
	}
	if !strings.Contains(out, `"display_name"`) {
		t.Error("JSON output should contain 'display_name' key")
	}
}
