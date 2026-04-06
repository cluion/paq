package provider

import (
	"fmt"
	"testing"
)

func TestFlatpakListPackages(t *testing.T) {
	output := "org.gimp.GIMP\t2.10.34\tGNU Image Manipulation Program\norg.mozilla.Firefox\t123.0\tWeb Browser\n"

	runner := func(name string, args ...string) (string, error) {
		return output, nil
	}

	p := &FlatpakProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pkgs) != 2 {
		t.Fatalf("expected 2 packages, got %d", len(pkgs))
	}

	if pkgs[0].Name != "org.gimp.GIMP" {
		t.Errorf("expected name 'org.gimp.GIMP', got %q", pkgs[0].Name)
	}
	if pkgs[0].Version != "2.10.34" {
		t.Errorf("expected version '2.10.34', got %q", pkgs[0].Version)
	}
	if pkgs[0].Desc != "GNU Image Manipulation Program" {
		t.Errorf("expected desc 'GNU Image Manipulation Program', got %q", pkgs[0].Desc)
	}
	if pkgs[1].Name != "org.mozilla.Firefox" {
		t.Errorf("expected name 'org.mozilla.Firefox', got %q", pkgs[1].Name)
	}
}

func TestFlatpakListOnlyApplicationName(t *testing.T) {
	output := "org.example.App\n"

	runner := func(name string, args ...string) (string, error) {
		return output, nil
	}

	p := &FlatpakProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pkgs) != 1 {
		t.Fatalf("expected 1 package, got %d", len(pkgs))
	}
	if pkgs[0].Name != "org.example.App" {
		t.Errorf("expected name 'org.example.App', got %q", pkgs[0].Name)
	}
	if pkgs[0].Version != "" {
		t.Errorf("expected empty version, got %q", pkgs[0].Version)
	}
	if pkgs[0].Desc != "" {
		t.Errorf("expected empty desc, got %q", pkgs[0].Desc)
	}
}

func TestFlatpakListSkipEmptyLines(t *testing.T) {
	output := "org.gimp.GIMP\t2.10.34\tImage Editor\n\n\n"

	runner := func(name string, args ...string) (string, error) {
		return output, nil
	}

	p := &FlatpakProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pkgs) != 1 {
		t.Fatalf("expected 1 package (empty lines skipped), got %d", len(pkgs))
	}
}

func TestFlatpakListEmpty(t *testing.T) {
	runner := func(name string, args ...string) (string, error) {
		return "", nil
	}

	p := &FlatpakProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pkgs) != 0 {
		t.Errorf("expected 0 packages, got %d", len(pkgs))
	}
}

func TestFlatpakListRunnerError(t *testing.T) {
	runner := func(name string, args ...string) (string, error) {
		return "", fmt.Errorf("flatpak not found")
	}

	p := &FlatpakProvider{runner: runner}
	_, err := p.List()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFlatpakName(t *testing.T) {
	p := &FlatpakProvider{runner: defaultRunner}
	if p.Name() != "flatpak" {
		t.Errorf("expected 'flatpak', got %q", p.Name())
	}
}

func TestFlatpakDisplayName(t *testing.T) {
	p := &FlatpakProvider{runner: defaultRunner}
	if p.DisplayName() != "Flatpak" {
		t.Errorf("expected 'Flatpak', got %q", p.DisplayName())
	}
}
