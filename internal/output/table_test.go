package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cluion/paq/internal/provider"
)

func TestTableFormatter_FormatPackages(t *testing.T) {
	pkgs := []provider.Package{
		{Name: "git", Version: "2.40.0", Type: "formula"},
		{Name: "firefox", Version: "146.0.1", Type: "cask"},
	}

	var buf bytes.Buffer
	f := &TableFormatter{}
	err := f.FormatPackages(&buf, "brew", pkgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()

	if !strings.Contains(out, "git") {
		t.Error("output should contain package name 'git'")
	}
	if !strings.Contains(out, "2.40.0") {
		t.Error("output should contain version '2.40.0'")
	}
	if !strings.Contains(out, "formula") {
		t.Error("output should contain type 'formula'")
	}
	if !strings.Contains(out, "cask") {
		t.Error("output should contain type 'cask'")
	}
	if !strings.Contains(out, "Total: 2") {
		t.Error("output should contain 'Total: 2'")
	}
}

func TestTableFormatter_FormatPackagesEmpty(t *testing.T) {
	var buf bytes.Buffer
	f := &TableFormatter{}
	err := f.FormatPackages(&buf, "brew", []provider.Package{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Total: 0") {
		t.Error("output should contain 'Total: 0'")
	}
}

func TestTableFormatter_FormatPackagesNoType(t *testing.T) {
	pkgs := []provider.Package{
		{Name: "curl", Version: "8.1.0"},
	}

	var buf bytes.Buffer
	f := &TableFormatter{}
	err := f.FormatPackages(&buf, "npm", pkgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "curl") {
		t.Error("output should contain package name 'curl'")
	}
	if strings.Contains(out, "Type") {
		t.Error("output should not contain Type column when no packages have type")
	}
}

func TestTableFormatter_FormatSources(t *testing.T) {
	sources := []SourceInfo{
		{Name: "brew", DisplayName: "Homebrew", Available: true},
		{Name: "snap", DisplayName: "Snap", Available: false},
	}

	var buf bytes.Buffer
	f := &TableFormatter{}
	err := f.FormatSources(&buf, sources)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()

	if !strings.Contains(out, "Homebrew") {
		t.Error("output should contain display name 'Homebrew'")
	}
	if !strings.Contains(out, "available") {
		t.Error("output should contain 'available'")
	}
	if !strings.Contains(out, "not found") {
		t.Error("output should contain 'not found'")
	}
	if !strings.Contains(out, "Total: 2") {
		t.Error("output should contain 'Total: 2'")
	}
}

func TestTableFormatter_FormatSourcesEmpty(t *testing.T) {
	var buf bytes.Buffer
	f := &TableFormatter{}
	err := f.FormatSources(&buf, []SourceInfo{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Total: 0") {
		t.Error("output should contain 'Total: 0'")
	}
}
