package provider

import (
	"fmt"
	"testing"
)

func TestAptListPackages(t *testing.T) {
	output := "curl\t8.12.1-1\tcommand line tool for transferring data with URL syntax\n" +
		"git\t1:2.47.0-1\tfast, scalable, distributed revision control system\n"

	runner := func(name string, args ...string) (string, error) {
		return output, nil
	}

	p := &AptProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pkgs) != 2 {
		t.Fatalf("expected 2 packages, got %d", len(pkgs))
	}

	if pkgs[0].Name != "curl" {
		t.Errorf("expected name 'curl', got %q", pkgs[0].Name)
	}
	if pkgs[0].Version != "8.12.1-1" {
		t.Errorf("expected version '8.12.1-1', got %q", pkgs[0].Version)
	}
	if pkgs[0].Desc != "command line tool for transferring data with URL syntax" {
		t.Errorf("unexpected desc: %q", pkgs[0].Desc)
	}

	if pkgs[1].Name != "git" {
		t.Errorf("expected name 'git', got %q", pkgs[1].Name)
	}
}

func TestAptListNoDesc(t *testing.T) {
	output := "base-files\t12.4\n"

	runner := func(name string, args ...string) (string, error) {
		return output, nil
	}

	p := &AptProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pkgs) != 1 {
		t.Fatalf("expected 1 package, got %d", len(pkgs))
	}
	if pkgs[0].Desc != "" {
		t.Errorf("expected empty desc, got %q", pkgs[0].Desc)
	}
}

func TestAptListEmpty(t *testing.T) {
	runner := func(name string, args ...string) (string, error) {
		return "", nil
	}

	p := &AptProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pkgs) != 0 {
		t.Errorf("expected 0 packages, got %d", len(pkgs))
	}
}

func TestAptListRunnerError(t *testing.T) {
	runner := func(name string, args ...string) (string, error) {
		return "", fmt.Errorf("dpkg-query not found")
	}

	p := &AptProvider{runner: runner}
	_, err := p.List()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAptListTrailingNewlines(t *testing.T) {
	output := "vim\t9.1.0\tVi IMproved\n\n\n"

	runner := func(name string, args ...string) (string, error) {
		return output, nil
	}

	p := &AptProvider{runner: runner}
	pkgs, err := p.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pkgs) != 1 {
		t.Fatalf("expected 1 package, got %d", len(pkgs))
	}
	if pkgs[0].Name != "vim" {
		t.Errorf("expected name 'vim', got %q", pkgs[0].Name)
	}
}

func TestAptName(t *testing.T) {
	p := &AptProvider{runner: defaultRunner}
	if p.Name() != "apt" {
		t.Errorf("expected 'apt', got %q", p.Name())
	}
}

func TestAptDisplayName(t *testing.T) {
	p := &AptProvider{runner: defaultRunner}
	if p.DisplayName() != "APT/dpkg" {
		t.Errorf("expected 'APT/dpkg', got %q", p.DisplayName())
	}
}
