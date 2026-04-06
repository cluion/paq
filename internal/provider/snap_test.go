package provider

import (
	"fmt"
	"testing"
)

func TestSnapListPackages(t *testing.T) {
	output := `Name  Version  Rev  Tracking       Publisher  Notes
core  16-2.58  155  latest/stable  canonical  core
curl  7.87    111  latest/stable  curl       -
`
	runner := func(name string, args ...string) (string, error) {
		return output, nil
	}

	p := &SnapProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pkgs) != 2 {
		t.Fatalf("expected 2 packages, got %d", len(pkgs))
	}

	if pkgs[0].Name != "core" {
		t.Errorf("expected name 'core', got %q", pkgs[0].Name)
	}
	if pkgs[0].Version != "16-2.58" {
		t.Errorf("expected version '16-2.58', got %q", pkgs[0].Version)
	}
	if pkgs[1].Name != "curl" {
		t.Errorf("expected name 'curl', got %q", pkgs[1].Name)
	}
	if pkgs[1].Version != "7.87" {
		t.Errorf("expected version '7.87', got %q", pkgs[1].Version)
	}
}

func TestSnapListSkipMalformedLines(t *testing.T) {
	output := `Name  Version  Rev
core  16-2.58  155
badline
curl  7.87  111
`
	runner := func(name string, args ...string) (string, error) {
		return output, nil
	}

	p := &SnapProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pkgs) != 2 {
		t.Fatalf("expected 2 packages (malformed line skipped), got %d", len(pkgs))
	}
}

func TestSnapListEmpty(t *testing.T) {
	output := `Name  Version  Rev  Tracking  Publisher  Notes
`
	runner := func(name string, args ...string) (string, error) {
		return output, nil
	}

	p := &SnapProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pkgs) != 0 {
		t.Errorf("expected 0 packages, got %d", len(pkgs))
	}
}

func TestSnapListRunnerError(t *testing.T) {
	runner := func(name string, args ...string) (string, error) {
		return "", fmt.Errorf("snap not found")
	}

	p := &SnapProvider{runner: runner}
	_, err := p.List()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSnapName(t *testing.T) {
	p := &SnapProvider{runner: defaultRunner}
	if p.Name() != "snap" {
		t.Errorf("expected 'snap', got %q", p.Name())
	}
}

func TestSnapDisplayName(t *testing.T) {
	p := &SnapProvider{runner: defaultRunner}
	if p.DisplayName() != "Snap" {
		t.Errorf("expected 'Snap', got %q", p.DisplayName())
	}
}
