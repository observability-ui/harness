package cli

import (
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd(version string) *cobra.Command {
	root := &cobra.Command{
		Use:          "obs",
		Short:        "Observability UI development tool",
		Long:         "A CLI tool to run recipes for developing, deploying, and managing Observability UI projects.",
		SilenceUsage: true,
	}

	root.AddCommand(newVersionCmd(version))
	root.AddCommand(newListCmd())
	root.AddCommand(newMixerCmd("start", "Start development processes"))
	root.AddCommand(newMixerCmd("deploy", "Deploy to an OpenShift cluster"))
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
