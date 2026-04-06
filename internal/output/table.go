package output

import (
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/cluion/paq/internal/provider"
)

// TableFormatter renders output as borderless tables.
type TableFormatter struct{}

// FormatPackages writes packages as a borderless table with a total count footer.
func (t *TableFormatter) FormatPackages(w io.Writer, source string, pkgs []provider.Package) error {
	tbl := table.New().
		Headers("Name", "Version", "Description").
		BorderLeft(false).
		BorderRight(false).
		BorderTop(false).
		BorderBottom(false).
		BorderHeader(false).
		BorderRow(false).
		BorderColumn(false)

	for _, p := range pkgs {
		tbl.Row(p.Name, p.Version, p.Desc)
	}

	if _, err := fmt.Fprintln(w, tbl.Render()); err != nil {
		return err
	}
	_, _ = fmt.Fprintf(w, "Total: %d\n", len(pkgs))
	return nil
}

// FormatSources writes source info as a borderless table with availability status.
func (t *TableFormatter) FormatSources(w io.Writer, sources []SourceInfo) error {
	tbl := table.New().
		Headers("Source", "Display Name", "Status").
		BorderLeft(false).
		BorderRight(false).
		BorderTop(false).
		BorderBottom(false).
		BorderHeader(false).
		BorderRow(false).
		BorderColumn(false)

	for _, s := range sources {
		status := "not found"
		if s.Available {
			status = "available"
		}
		tbl.Row(s.Name, s.DisplayName, status)
	}

	if _, err := fmt.Fprintln(w, tbl.Render()); err != nil {
		return err
	}
	_, _ = fmt.Fprintf(w, "Total: %d\n", len(sources))
	return nil
}
