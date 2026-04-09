package provider

import (
	"encoding/json"
)

// BrewProvider queries Homebrew for installed formulae and casks.
type BrewProvider struct {
	runner CommandRunner
}

func init() {
	Register(&BrewProvider{runner: defaultRunner})
}

func (b *BrewProvider) Name() string        { return "brew" }
func (b *BrewProvider) DisplayName() string  { return "Homebrew" }
func (b *BrewProvider) Detect() bool         { return commandExists("brew") }

func (b *BrewProvider) List() ([]Package, error) {
	out, err := b.runner("brew", "info", "--installed", "--json=v2")
	if err != nil {
		return nil, err
	}

	var data struct {
		Formulae []struct {
			Name      string `json:"name"`
			Installed []struct {
				Version string `json:"version"`
			} `json:"installed"`
			Desc string `json:"desc"`
		} `json:"formulae"`
		Casks []struct {
			Name      string `json:"token"`
			Installed string `json:"installed"`
			Desc      string `json:"desc"`
		} `json:"casks"`
	}

	if err := json.Unmarshal([]byte(out), &data); err != nil {
		return nil, err
	}

	pkgs := make([]Package, 0, len(data.Formulae)+len(data.Casks))

	for _, f := range data.Formulae {
		version := ""
		if len(f.Installed) > 0 {
			version = f.Installed[0].Version
		}
		pkgs = append(pkgs, Package{
			Name:    f.Name,
			Version: version,
			Desc:    f.Desc,
			Type:    "formula",
		})
	}

	for _, c := range data.Casks {
		pkgs = append(pkgs, Package{
			Name:    c.Name,
			Version: c.Installed,
			Desc:    c.Desc,
			Type:    "cask",
		})
	}

	return pkgs, nil
}
