package output

import (
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/cluion/paq/internal/provider"
)

const (
	maxNameLen    = 20
	maxVersionLen = 14
)

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

var (
	baseStyle      = lipgloss.NewStyle().Padding(0, 1)
	headerStyle    = lipgloss.NewStyle().Bold(true).Padding(0, 1)
	nameStyle      = lipgloss.NewStyle().Padding(0, 1)
	versionStyle   = lipgloss.NewStyle().Padding(0, 1)
	typeStyle      = lipgloss.NewStyle().Padding(0, 1)
	availableStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Padding(0, 1)  // green
	notFoundStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("246")).Padding(0, 1) // gray
)

// TableFormatter renders output as styled tables.
type TableFormatter struct{}

// FormatPackages writes packages as a styled table with a total count footer.
func (t *TableFormatter) FormatPackages(w io.Writer, source string, pkgs []provider.Package) error {
	hasType := false
	for _, p := range pkgs {
		if p.Type != "" {
			hasType = true
			break
		}
	}

	headers := []string{"Name", "Version"}
	if hasType {
		headers = append(headers, "Type")
	}

	tbl := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		BorderHeader(true).
		BorderColumn(true).
		BorderRow(false).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}
			switch col {
			case 0:
				return nameStyle
			case 1:
				return versionStyle
			default:
				return typeStyle
			}
		}).
		Headers(headers...)

	for _, p := range pkgs {
		row := []string{truncate(p.Name, maxNameLen), truncate(p.Version, maxVersionLen)}
		if hasType {
			row = append(row, p.Type)
		}
		tbl.Row(row...)
	}

	if _, err := fmt.Fprintln(w, tbl.Render()); err != nil {
		return err
	}
	_, _ = fmt.Fprintf(w, "Total: %d\n", len(pkgs))
	return nil
}

// FormatSources writes source info as a styled table with availability status.
func (t *TableFormatter) FormatSources(w io.Writer, sources []SourceInfo) error {
	tbl := table.New().
		Border(lipgloss.RoundedBorder()).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}
			return baseStyle
		}).
		Headers("Source", "Display Name", "Status")

	for _, s := range sources {
		var status string
		if s.Available {
			status = availableStyle.Render("available")
		} else {
			status = notFoundStyle.Render("not found")
		}
		tbl.Row(s.Name, s.DisplayName, status)
	}

	if _, err := fmt.Fprintln(w, tbl.Render()); err != nil {
		return err
	}
	_, _ = fmt.Fprintf(w, "Total: %d\n", len(sources))
	return nil
}
