package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"obs/internal/process"
	"obs/internal/state"

	"github.com/spf13/cobra"
)

func newCleanupCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "cleanup [recipe]",
		Short: "Stop all processes and clean up state",
		RunE: func(cmd *cobra.Command, args []string) error {
			store := state.NewStore(state.DefaultStateDir())
			rs, err := store.Load()
			if err != nil {
				fmt.Println("Nothing to clean up.")
				return nil
			}
			rs = store.FilterAlive(rs)

			if len(rs.Processes) > 0 && !force {
				fmt.Printf("This will stop %d running process(es). Continue? [y/N] ", len(rs.Processes))
				reader := bufio.NewReader(os.Stdin)
				answer, _ := reader.ReadString('\n')
				if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(answer)), "y") {
					fmt.Println("Aborted.")
					return nil
				}
			}

			for _, p := range rs.Processes {
				if process.IsAlive(p.PID) {
					fmt.Printf("Stopping %s (PID %d)…\n", p.Name, p.PID)
					process.KillGroup(p.PID)
				}
			}
			store.Clean()
			fmt.Println("Cleanup complete.")
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "skip confirmation")
	return cmd
}
