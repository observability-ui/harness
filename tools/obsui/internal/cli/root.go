package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	nonInteractive bool
	outputJSON     bool
	dryRun         bool
	detach         bool
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
	root.PersistentFlags().BoolVar(&detach, "detach", false, "start processes in background and exit")

	root.AddCommand(newVersionCmd(version))
	root.AddCommand(newListCmd())
	root.AddCommand(newRecipeCmd("start", "Start development processes"))
	root.AddCommand(newRecipeCmd("deploy", "Deploy to an OpenShift cluster"))
	root.AddCommand(newStatusCmd())
	root.AddCommand(newAttachCmd())
	root.AddCommand(newCleanupCmd())

	return root
}

func Execute(version string) {
	root := NewRootCmd(version)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
