package output

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/cluion/paq/internal/provider"
)

func TestJSONFormatter_FormatPackages(t *testing.T) {
	pkgs := []provider.Package{
		{Name: "git", Version: "2.40.0", Desc: "Distributed version control system"},
		{Name: "node", Version: "20.1.0", Desc: "Platform built on V8"},
	}

	var buf bytes.Buffer
	f := &JSONFormatter{}
	err := f.FormatPackages(&buf, "npm", pkgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if string(result["source"]) != `"npm"` {
		t.Errorf("expected source \"npm\", got %s", result["source"])
	}
	if string(result["count"]) != "2" {
		t.Errorf("expected count 2, got %s", result["count"])
	}
	if result["packages"] == nil {
		t.Error("expected packages field")
	}
}

func TestJSONFormatter_FormatPackagesEmpty(t *testing.T) {
	var buf bytes.Buffer
	f := &JSONFormatter{}
	err := f.FormatPackages(&buf, "brew", []provider.Package{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if string(result["count"]) != "0" {
		t.Errorf("expected count 0, got %s", result["count"])
	}
}

func TestJSONFormatter_FormatPackagesOmitEmptyDesc(t *testing.T) {
	pkgs := []provider.Package{
		{Name: "curl", Version: "8.1.0"},
	}

	var buf bytes.Buffer
	f := &JSONFormatter{}
	err := f.FormatPackages(&buf, "brew", pkgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte(`"name"`)) {
		t.Error("output should contain package name")
	}
	if bytes.Contains(buf.Bytes(), []byte(`"desc"`)) {
		t.Error("empty desc should be omitted")
	}
}

func TestJSONFormatter_FormatSources(t *testing.T) {
	sources := []SourceInfo{
		{Name: "brew", DisplayName: "Homebrew", Available: true},
		{Name: "snap", DisplayName: "Snap", Available: false},
	}

	var buf bytes.Buffer
	f := &JSONFormatter{}
	err := f.FormatSources(&buf, sources)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result struct {
		Sources []SourceInfo `json:"sources"`
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if len(result.Sources) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(result.Sources))
	}
	if result.Sources[0].Name != "brew" {
		t.Errorf("expected first source name \"brew\", got %q", result.Sources[0].Name)
	}
	if !result.Sources[0].Available {
		t.Error("brew should be available")
	}
	if result.Sources[1].Available {
		t.Error("snap should not be available")
	}
}

func TestJSONFormatter_FormatSourcesEmpty(t *testing.T) {
	var buf bytes.Buffer
	f := &JSONFormatter{}
	err := f.FormatSources(&buf, []SourceInfo{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result struct {
		Sources []SourceInfo `json:"sources"`
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if len(result.Sources) != 0 {
		t.Errorf("expected 0 sources, got %d", len(result.Sources))
	}
}

func TestNewFormatter(t *testing.T) {
	if _, ok := New(false).(*TableFormatter); !ok {
		t.Error("New(false) should return *TableFormatter")
	}
	if _, ok := New(true).(*JSONFormatter); !ok {
		t.Error("New(true) should return *JSONFormatter")
	}
}
