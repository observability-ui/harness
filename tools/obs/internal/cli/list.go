package cli

import (
	"fmt"
	"strings"

	"obs/internal/mixer"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available recipes",
		Run: func(cmd *cobra.Command, args []string) {
			entries := mixer.ListAllRecipes()
			if len(entries) == 0 {
				fmt.Println("No recipes registered.")
				return
			}
			for _, entry := range entries {
				aliases := ""
				if len(entry.Aliases) > 0 {
					aliases = " (" + strings.Join(entry.Aliases, ", ") + ")"
				}
				components := strings.Join(entry.Components, ", ")
				fmt.Printf("  %-8s %-25s%s  [%s]\n", entry.Command, entry.Name, aliases, components)
			}
		},
	}
}
