package cli

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/cluion/paq/internal/provider"
)

type queryTestProvider struct {
	name        string
	displayName string
	detect      bool
	pkgs        []provider.Package
	listErr     error
}

func (p *queryTestProvider) Name() string                    { return p.name }
func (p *queryTestProvider) DisplayName() string              { return p.displayName }
func (p *queryTestProvider) Detect() bool                     { return p.detect }
func (p *queryTestProvider) List() ([]provider.Package, error) { return p.pkgs, p.listErr }

func TestProviderCommandAvailable(t *testing.T) {
	provider.Register(&queryTestProvider{
		name:        "testquery-available",
		displayName: "TestAvailable",
		detect:      true,
		pkgs: []provider.Package{
			{Name: "pkg1", Version: "1.0.0"},
		},
	})

	RegisterProviderCommands()

	jsonOutput = false
	out, err := executeCommand(rootCmd, "testquery-available")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "pkg1") {
		t.Error("output should contain 'pkg1'")
	}
	if !strings.Contains(out, "Total: 1") {
		t.Error("output should contain 'Total: 1'")
	}
}

func TestProviderCommandNotAvailable(t *testing.T) {
	provider.Register(&queryTestProvider{
		name:        "testquery-unavail",
		displayName: "TestUnavail",
		detect:      false,
	})

	RegisterProviderCommands()

	_, err := executeCommand(rootCmd, "testquery-unavail")
	if err == nil {
		t.Fatal("expected error when provider is not available")
	}

	expected := fmt.Sprintf("%s 未安裝或不在 PATH 中", "TestUnavail")
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("error should contain %q, got %q", expected, err.Error())
	}
}

func TestProviderCommandJSON(t *testing.T) {
	provider.Register(&queryTestProvider{
		name:        "testquery-json",
		displayName: "TestJSON",
		detect:      true,
		pkgs: []provider.Package{
			{Name: "pkg2", Version: "2.0.0", Desc: "A test package"},
		},
	})

	RegisterProviderCommands()

	jsonOutput = true
	rootCmd.SetArgs([]string{"testquery-json"})
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	err := rootCmd.Execute()
	jsonOutput = false

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"source"`) {
		t.Error("JSON output should contain 'source' key")
	}
	if !strings.Contains(out, `"packages"`) {
		t.Error("JSON output should contain 'packages' key")
	}
	if !strings.Contains(out, `"pkg2"`) {
		t.Error("JSON output should contain package name 'pkg2'")
	}
}

func TestProviderCommandListError(t *testing.T) {
	provider.Register(&queryTestProvider{
		name:        "testquery-err",
		displayName: "TestErr",
		detect:      true,
		listErr:     fmt.Errorf("something went wrong"),
	})

	RegisterProviderCommands()

	_, err := executeCommand(rootCmd, "testquery-err")
	if err == nil {
		t.Fatal("expected error from List()")
	}
	if !strings.Contains(err.Error(), "something went wrong") {
		t.Errorf("error should contain original message, got %q", err.Error())
	}
}
