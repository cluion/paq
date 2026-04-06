package provider

import (
	"fmt"
	"sort"
	"testing"
)

func TestNpmListPackages(t *testing.T) {
	jsonOut := `{
		"dependencies": {
			"typescript": {"version": "5.0.4"},
			"eslint": {"version": "8.40.0"},
			"prettier": {"version": "2.8.8"}
		}
	}`

	runner := func(name string, args ...string) (string, error) {
		return jsonOut, nil
	}

	p := &NpmProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pkgs) != 3 {
		t.Fatalf("expected 3 packages, got %d", len(pkgs))
	}

	// Sort for deterministic assertion
	sort.Slice(pkgs, func(i, j int) bool { return pkgs[i].Name < pkgs[j].Name })

	if pkgs[0].Name != "eslint" {
		t.Errorf("expected 'eslint', got %q", pkgs[0].Name)
	}
	if pkgs[0].Version != "8.40.0" {
		t.Errorf("expected version '8.40.0', got %q", pkgs[0].Version)
	}
	if pkgs[1].Name != "prettier" {
		t.Errorf("expected 'prettier', got %q", pkgs[1].Name)
	}
	if pkgs[2].Name != "typescript" {
		t.Errorf("expected 'typescript', got %q", pkgs[2].Name)
	}
}

func TestNpmListNonZeroExitWithOutput(t *testing.T) {
	jsonOut := `{
		"dependencies": {
			"npm": {"version": "9.6.7"}
		}
	}`

	runner := func(name string, args ...string) (string, error) {
		return jsonOut, fmt.Errorf("exit code 1")
	}

	p := &NpmProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("should tolerate non-zero exit when stdout has data: %v", err)
	}

	if len(pkgs) != 1 {
		t.Fatalf("expected 1 package, got %d", len(pkgs))
	}
	if pkgs[0].Name != "npm" {
		t.Errorf("expected 'npm', got %q", pkgs[0].Name)
	}
}

func TestNpmListNonZeroExitEmptyStdout(t *testing.T) {
	runner := func(name string, args ...string) (string, error) {
		return "", fmt.Errorf("exit code 1")
	}

	p := &NpmProvider{runner: runner}
	_, err := p.List()
	if err == nil {
		t.Fatal("expected error when stdout is empty and exit code is non-zero")
	}
}

func TestNpmListEmpty(t *testing.T) {
	runner := func(name string, args ...string) (string, error) {
		return `{"dependencies": {}}`, nil
	}

	p := &NpmProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pkgs) != 0 {
		t.Errorf("expected 0 packages, got %d", len(pkgs))
	}
}

func TestNpmListInvalidJSON(t *testing.T) {
	runner := func(name string, args ...string) (string, error) {
		return "not json", nil
	}

	p := &NpmProvider{runner: runner}
	_, err := p.List()
	if err == nil {
		t.Fatal("expected JSON parse error")
	}
}

func TestNpmName(t *testing.T) {
	p := &NpmProvider{runner: defaultRunner}
	if p.Name() != "npm" {
		t.Errorf("expected 'npm', got %q", p.Name())
	}
}

func TestNpmDisplayName(t *testing.T) {
	p := &NpmProvider{runner: defaultRunner}
	if p.DisplayName() != "npm" {
		t.Errorf("expected 'npm', got %q", p.DisplayName())
	}
}
