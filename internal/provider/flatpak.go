package provider

import (
	"strings"
)

// FlatpakProvider queries flatpak for installed applications.
type FlatpakProvider struct {
	runner CommandRunner
}

func init() {
	Register(&FlatpakProvider{runner: defaultRunner})
}

func (f *FlatpakProvider) Name() string        { return "flatpak" }
func (f *FlatpakProvider) DisplayName() string  { return "Flatpak" }
func (f *FlatpakProvider) Detect() bool         { return commandExists("flatpak") }

func (f *FlatpakProvider) List() ([]Package, error) {
	out, err := f.runner("flatpak", "list", "--app", "--columns=application,version,description")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	pkgs := make([]Package, 0, len(lines))

	for _, line := range lines {
		parts := strings.SplitN(line, "\t", 3)

		name := strings.TrimSpace(parts[0])
		if name == "" {
			continue
		}

		version := ""
		if len(parts) > 1 {
			version = strings.TrimSpace(parts[1])
		}

		desc := ""
		if len(parts) > 2 {
			desc = strings.TrimSpace(parts[2])
		}

		pkgs = append(pkgs, Package{
			Name:    name,
			Version: version,
			Desc:    desc,
		})
	}

	return pkgs, nil
}
