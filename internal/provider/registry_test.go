package provider

import (
	"testing"
)

type mockProvider struct {
	name        string
	displayName string
	detect      bool
	pkgs        []Package
	listErr     error
}

func (m *mockProvider) Name() string                { return m.name }
func (m *mockProvider) DisplayName() string          { return m.displayName }
func (m *mockProvider) Detect() bool                 { return m.detect }
func (m *mockProvider) List() ([]Package, error)     { return m.pkgs, m.listErr }

func newRegistry() *Registry {
	return &Registry{providers: make(map[string]Provider)}
}

func TestRegisterAndGet(t *testing.T) {
	r := newRegistry()
	p := &mockProvider{name: "test", displayName: "Test"}

	r.register(p)

	got, ok := r.get("test")
	if !ok {
		t.Fatal("expected provider to be found")
	}
	if got.Name() != "test" {
		t.Errorf("expected name %q, got %q", "test", got.Name())
	}
}

func TestGetNotFound(t *testing.T) {
	r := newRegistry()

	_, ok := r.get("nonexistent")
	if ok {
		t.Error("expected provider not to be found")
	}
}

func TestRegisterDuplicatePanics(t *testing.T) {
	r := newRegistry()
	p := &mockProvider{name: "dup", displayName: "Dup"}
	r.register(p)

	defer func() {
		rec := recover()
		if rec == nil {
			t.Fatal("expected panic on duplicate registration")
		}
		msg, ok := rec.(string)
		if !ok {
			t.Fatalf("expected string panic, got %T: %v", rec, rec)
		}
		if msg != "provider already registered: dup" {
			t.Errorf("unexpected panic message: %s", msg)
		}
	}()

	r.register(&mockProvider{name: "dup", displayName: "Dup2"})
}

func TestAllSorted(t *testing.T) {
	r := newRegistry()
	r.register(&mockProvider{name: "snap", displayName: "Snap"})
	r.register(&mockProvider{name: "brew", displayName: "Homebrew"})
	r.register(&mockProvider{name: "npm", displayName: "npm"})

	all := r.all()

	if len(all) != 3 {
		t.Fatalf("expected 3 providers, got %d", len(all))
	}
	names := []string{all[0].Name(), all[1].Name(), all[2].Name()}
	expected := []string{"brew", "npm", "snap"}
	for i, exp := range expected {
		if names[i] != exp {
			t.Errorf("position %d: expected %q, got %q", i, exp, names[i])
		}
	}
}

func TestAvailable(t *testing.T) {
	r := newRegistry()
	r.register(&mockProvider{name: "brew", displayName: "Homebrew", detect: true})
	r.register(&mockProvider{name: "snap", displayName: "Snap", detect: false})
	r.register(&mockProvider{name: "npm", displayName: "npm", detect: true})

	avail := r.available()

	if len(avail) != 2 {
		t.Fatalf("expected 2 available providers, got %d", len(avail))
	}
	for _, p := range avail {
		if p.Name() == "snap" {
			t.Error("snap should not be available")
		}
	}
}

func TestAllEmpty(t *testing.T) {
	r := newRegistry()

	all := r.all()
	if len(all) != 0 {
		t.Errorf("expected empty slice, got %d", len(all))
	}
}
