package provider

import (
	"fmt"
	"testing"
)

func fakeBrewJSON(formulae int, casks int) string {
	formulaeJSON := ""
	for i := 0; i < formulae; i++ {
		formulaeJSON += fmt.Sprintf(`{
			"name": "formula-%d",
			"installed": [{"version": "%d.0.0"}],
			"desc": "Formula %d description"
		}`, i, i, i)
		if i < formulae-1 {
			formulaeJSON += ","
		}
	}

	casksJSON := ""
	for i := 0; i < casks; i++ {
		casksJSON += fmt.Sprintf(`{
			"token": "cask-%d",
			"installed": "%d.0.0",
			"desc": "Cask %d description"
		}`, i, i, i)
		if i < casks-1 {
			casksJSON += ","
		}
	}

	return fmt.Sprintf(`{
		"formulae": [%s],
		"casks": [%s]
	}`, formulaeJSON, casksJSON)
}

func TestBrewListFormulaeAndCasks(t *testing.T) {
	jsonOut := fakeBrewJSON(1, 1)
	runner := func(name string, args ...string) (string, error) {
		return jsonOut, nil
	}

	p := &BrewProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pkgs) != 2 {
		t.Fatalf("expected 2 packages, got %d", len(pkgs))
	}

	formula := pkgs[0]
	if formula.Name != "formula-0" {
		t.Errorf("expected name 'formula-0', got %q", formula.Name)
	}
	if formula.Version != "0.0.0" {
		t.Errorf("expected version '0.0.0', got %q", formula.Version)
	}
	if formula.Desc != "Formula 0 description" {
		t.Errorf("expected desc 'Formula 0 description', got %q", formula.Desc)
	}

	cask := pkgs[1]
	if cask.Name != "cask-0" {
		t.Errorf("expected name 'cask-0', got %q", cask.Name)
	}
	if cask.Version != "0.0.0" {
		t.Errorf("expected version '0.0.0', got %q", cask.Version)
	}
}

func TestBrewListEmptyInstalled(t *testing.T) {
	jsonOut := `{
		"formulae": [{"name": "autoconf", "installed": [], "desc": "Automatic configure script builder"}],
		"casks": []
	}`

	runner := func(name string, args ...string) (string, error) {
		return jsonOut, nil
	}

	p := &BrewProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pkgs) != 1 {
		t.Fatalf("expected 1 package, got %d", len(pkgs))
	}
	if pkgs[0].Version != "" {
		t.Errorf("expected empty version, got %q", pkgs[0].Version)
	}
}

func TestBrewListEmpty(t *testing.T) {
	runner := func(name string, args ...string) (string, error) {
		return `{"formulae":[],"casks":[]}`, nil
	}

	p := &BrewProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pkgs) != 0 {
		t.Errorf("expected 0 packages, got %d", len(pkgs))
	}
}

func TestBrewListRunnerError(t *testing.T) {
	runner := func(name string, args ...string) (string, error) {
		return "", fmt.Errorf("brew not found")
	}

	p := &BrewProvider{runner: runner}
	_, err := p.List()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBrewListInvalidJSON(t *testing.T) {
	runner := func(name string, args ...string) (string, error) {
		return "not json", nil
	}

	p := &BrewProvider{runner: runner}
	_, err := p.List()
	if err == nil {
		t.Fatal("expected JSON parse error")
	}
}

func TestBrewName(t *testing.T) {
	p := &BrewProvider{runner: defaultRunner}
	if p.Name() != "brew" {
		t.Errorf("expected 'brew', got %q", p.Name())
	}
}

func TestBrewDisplayName(t *testing.T) {
	p := &BrewProvider{runner: defaultRunner}
	if p.DisplayName() != "Homebrew" {
		t.Errorf("expected 'Homebrew', got %q", p.DisplayName())
	}
}
