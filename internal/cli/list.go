package cli

import (
	"github.com/cluion/paq/internal/output"
	"github.com/cluion/paq/internal/provider"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered package sources",
	RunE: func(cmd *cobra.Command, args []string) error {
		all := provider.All()
		sources := make([]output.SourceInfo, 0, len(all))
		for _, p := range all {
			sources = append(sources, output.SourceInfo{
				Name:        p.Name(),
				DisplayName: p.DisplayName(),
				Available:   p.Detect(),
			})
		}

		formatter := output.New(GetJSONOutput())
		return formatter.FormatSources(cmd.OutOrStdout(), sources)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
