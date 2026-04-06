package output

import (
	"encoding/json"
	"io"

	"github.com/cluion/paq/internal/provider"
)

// JSONFormatter renders output as pretty-printed JSON.
type JSONFormatter struct{}

type packageEnvelope struct {
	Source   string            `json:"source"`
	Count    int               `json:"count"`
	Packages []provider.Package `json:"packages"`
}

type sourcesEnvelope struct {
	Sources []SourceInfo `json:"sources"`
}

// FormatPackages writes packages as JSON with source, count, and packages envelope.
func (j *JSONFormatter) FormatPackages(w io.Writer, source string, pkgs []provider.Package) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(packageEnvelope{
		Source:   source,
		Count:    len(pkgs),
		Packages: pkgs,
	})
}

// FormatSources writes source info as JSON with a sources envelope.
func (j *JSONFormatter) FormatSources(w io.Writer, sources []SourceInfo) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(sourcesEnvelope{
		Sources: sources,
	})
}
