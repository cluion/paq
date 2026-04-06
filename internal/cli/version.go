package cli

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// buildInfo returns (version, commit, date) using ldflags first,
// falling back to runtime/debug.ReadBuildInfo() for go install builds.
func buildInfo() (string, string, string) {
	v, c, d := version, commit, date

	if v == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok {
			if info.Main.Version != "" && info.Main.Version != "(devel)" {
				v = info.Main.Version
			}
			for _, s := range info.Settings {
				switch s.Key {
				case "vcs.revision":
					if c == "none" {
						c = s.Value
						if len(c) > 8 {
							c = c[:8]
						}
					}
				case "vcs.time":
					if d == "unknown" {
						d = s.Value
					}
				}
			}
		}
	}

	return v, c, d
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",
	Run: func(cmd *cobra.Command, args []string) {
		v, c, d := buildInfo()
		fmt.Fprintf(cmd.OutOrStdout(), "paq %s (commit: %s, built: %s)\n", v, c, d)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
