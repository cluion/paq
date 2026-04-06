package cli

import (
	"fmt"

	"github.com/cluion/paq/internal/output"
	"github.com/cluion/paq/internal/provider"
	"github.com/spf13/cobra"
)

// reservedNames are built-in command names that providers cannot override.
var reservedNames = map[string]bool{
	"list":       true,
	"version":    true,
	"help":       true,
	"completion": true,
	"upgrade":    true,
}

// RegisterProviderCommands creates a subcommand for each registered Provider.
func RegisterProviderCommands() {
	for _, p := range provider.All() {
		if reservedNames[p.Name()] {
			continue
		}
		addProviderCmd(p)
	}
}

func addProviderCmd(p provider.Provider) {
	name := p.Name()
	displayName := p.DisplayName()

	cmd := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("List packages from %s", displayName),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !p.Detect() {
				return fmt.Errorf("%s 未安裝或不在 PATH 中", displayName)
			}

			pkgs, err := p.List()
			if err != nil {
				return err
			}

			formatter := output.New(GetJSONOutput())
			return formatter.FormatPackages(cmd.OutOrStdout(), name, pkgs)
		},
	}

	rootCmd.AddCommand(cmd)
}
