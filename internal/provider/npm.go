package provider

import (
	"encoding/json"
	"strings"
)

// NpmProvider queries npm for globally installed packages.
type NpmProvider struct {
	runner CommandRunner
}

func init() {
	Register(&NpmProvider{runner: defaultRunner})
}

func (n *NpmProvider) Name() string        { return "npm" }
func (n *NpmProvider) DisplayName() string  { return "npm" }
func (n *NpmProvider) Detect() bool         { return commandExists("npm") }

func (n *NpmProvider) List() ([]Package, error) {
	out, err := n.runner("npm", "list", "-g", "--json", "--depth=0")

	// npm may return non-zero exit due to peer dependency warnings.
	// Only fail if stdout is also empty.
	if err != nil && strings.TrimSpace(out) == "" {
		return nil, err
	}

	var data struct {
		Dependencies map[string]struct {
			Version string `json:"version"`
		} `json:"dependencies"`
	}

	if err := json.Unmarshal([]byte(out), &data); err != nil {
		return nil, err
	}

	pkgs := make([]Package, 0, len(data.Dependencies))
	for name, info := range data.Dependencies {
		pkgs = append(pkgs, Package{
			Name:    name,
			Version: info.Version,
		})
	}

	return pkgs, nil
}
