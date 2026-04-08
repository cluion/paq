package provider

import (
	"strings"
)

// AptProvider queries dpkg for installed packages on Debian/Ubuntu systems.
type AptProvider struct {
	runner CommandRunner
}

func init() {
	Register(&AptProvider{runner: defaultRunner})
}

func (a *AptProvider) Name() string        { return "apt" }
func (a *AptProvider) DisplayName() string  { return "APT/dpkg" }
func (a *AptProvider) Detect() bool         { return commandExists("dpkg") }

func (a *AptProvider) List() ([]Package, error) {
	out, err := a.runner("dpkg-query",
		"-f", "${Package}\t${Version}\t${binary:Summary}\n",
		"-W",
	)
	if err != nil {
		return nil, err
	}

	var pkgs []Package
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.SplitN(line, "\t", 3)
		pkg := Package{Name: fields[0]}
		if len(fields) > 1 {
			pkg.Version = fields[1]
		}
		if len(fields) > 2 {
			pkg.Desc = fields[2]
		}
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}
