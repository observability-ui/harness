package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"obs/internal/state"

	"github.com/spf13/cobra"
)

func newAttachCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "attach",
		Short: "Attach to running processes and show their output",
		RunE: func(cmd *cobra.Command, args []string) error {
			store := state.NewStore(state.DefaultStateDir())
			rs, err := store.Load()
			if err != nil {
				return fmt.Errorf("no active processes to attach to")
			}
			rs = store.FilterAlive(rs)
			if len(rs.Processes) == 0 {
				return fmt.Errorf("no running processes found")
			}

			fmt.Println("Attached to running processes. Press Ctrl+C to detach.")
			state.PrintStatus(rs)

			// Wait for interrupt
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT)

			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-sigCh:
					fmt.Println("\nDetached.")
					return nil
				case <-ticker.C:
					rs, _ = store.Load()
					if rs != nil {
						rs = store.FilterAlive(rs)
						if len(rs.Processes) == 0 {
							fmt.Println("All processes have exited.")
							return nil
						}
					}
				}
			}
		},
	}
}
