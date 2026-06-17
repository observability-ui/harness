package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"obsui/internal/recipe"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available recipes",
		Run: func(cmd *cobra.Command, args []string) {
			recipes := recipe.DefaultRegistry.ListAll()
			if len(recipes) == 0 {
				fmt.Println("No recipes registered.")
				return
			}
			for _, r := range recipes {
				aliases := ""
				if len(r.Aliases()) > 0 {
					aliases = " (" + strings.Join(r.Aliases(), ", ") + ")"
				}
				fmt.Printf("  %-8s %-25s %s%s\n", r.Command(), r.Name(), r.Description(), aliases)
			}
		},
	}
}
