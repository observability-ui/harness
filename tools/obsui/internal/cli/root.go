package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	nonInteractive bool
	outputJSON     bool
	dryRun         bool
)

func NewRootCmd(version string) *cobra.Command {
	root := &cobra.Command{
		Use:   "obsui",
		Short: "Observability UI development tool",
		Long:  "A CLI tool to run recipes for developing, deploying, and managing Observability UI projects.",
		SilenceUsage: true,
	}

	root.PersistentFlags().BoolVar(&nonInteractive, "non-interactive", false, "force non-interactive mode")
	root.PersistentFlags().BoolVar(&outputJSON, "output-json", false, "emit JSON events instead of terminal output")
	root.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "show what would run without executing")

	root.AddCommand(newVersionCmd(version))
	root.AddCommand(newListCmd())

	return root
}

func Execute(version string) {
	root := NewRootCmd(version)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
