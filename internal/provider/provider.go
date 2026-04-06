package provider

// Package represents an installed package.
type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Desc    string `json:"desc,omitempty"`
}

// Provider is the interface that all package source providers must implement.
type Provider interface {
	// Name returns the unique identifier for the provider (e.g. "brew").
	Name() string

	// DisplayName returns the human-readable name (e.g. "Homebrew").
	DisplayName() string

	// Detect checks whether the provider's command is available on the system.
	Detect() bool

	// List returns all packages managed by this provider.
	List() ([]Package, error)
}
