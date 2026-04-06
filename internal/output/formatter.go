package output

import (
	"io"

	"github.com/cluion/paq/internal/provider"
)

// SourceInfo describes a registered provider with its detection status.
type SourceInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Available   bool   `json:"available"`
}

// Formatter renders package or source information to a writer.
type Formatter interface {
	FormatPackages(w io.Writer, source string, pkgs []provider.Package) error
	FormatSources(w io.Writer, sources []SourceInfo) error
}

// New returns a JSONFormatter when jsonOutput is true, otherwise a TableFormatter.
func New(jsonOutput bool) Formatter {
	if jsonOutput {
		return &JSONFormatter{}
	}
	return &TableFormatter{}
}
