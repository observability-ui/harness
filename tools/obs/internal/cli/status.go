package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"obs/internal/state"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status [recipe]",
		Short: "Show status of running processes",
		RunE: func(cmd *cobra.Command, args []string) error {
			store := state.NewStore(state.DefaultStateDir())
			rs, err := store.Load()
			if err != nil {
				fmt.Println("No active processes.")
				return nil
			}
			rs = store.FilterAlive(rs)
			state.PrintStatus(rs)
			return nil
		},
	}
}
