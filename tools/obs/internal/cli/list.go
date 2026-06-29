package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"obs/internal/recipe"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available recipes",
		Run: func(cmd *cobra.Command, args []string) {
			entries := recipe.DefaultRegistry.ListAll()
			if len(entries) == 0 {
				fmt.Println("No recipes registered.")
				return
			}
			for _, entry := range entries {
				aliases := ""
				if len(entry.Recipe.Aliases()) > 0 {
					aliases = " (" + strings.Join(entry.Recipe.Aliases(), ", ") + ")"
				}
				fmt.Printf("  %-8s %-25s %s%s\n", entry.Command, entry.Recipe.Name(), entry.Recipe.Description(), aliases)
			}
		},
	}
}
