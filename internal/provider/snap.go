package provider

import (
	"strings"
)

// SnapProvider queries snap for installed packages.
type SnapProvider struct {
	runner CommandRunner
}

func init() {
	Register(&SnapProvider{runner: defaultRunner})
}

func (s *SnapProvider) Name() string        { return "snap" }
func (s *SnapProvider) DisplayName() string  { return "Snap" }
func (s *SnapProvider) Detect() bool         { return commandExists("snap") }

func (s *SnapProvider) List() ([]Package, error) {
	out, err := s.runner("snap", "list")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")

	// Skip header line
	if len(lines) <= 1 {
		return nil, nil
	}

	pkgs := make([]Package, 0, len(lines)-1)
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		pkgs = append(pkgs, Package{
			Name:    fields[0],
			Version: fields[1],
		})
	}

	return pkgs, nil
}
