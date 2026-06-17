package runner

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"obsui/internal/process"
	"obsui/internal/recipe"
	"obsui/internal/types"
	"obsui/internal/ui"
)

type InteractiveRunner struct {
	program *tea.Program
}

func NewInteractive() *InteractiveRunner {
	return &InteractiveRunner{}
}

func (r *InteractiveRunner) Run(ctx context.Context, mgr *process.Manager, steps []*recipe.Step, updates chan<- types.StepUpdate) error {
	internalUpdates := make(chan types.StepUpdate, 100)
	model := ui.NewModel(mgr, steps, internalUpdates)

	r.program = tea.NewProgram(model, tea.WithAltScreen())

	// Start processes in a goroutine and send updates
	go func() {
		defer close(internalUpdates)
		defer close(updates)

		for _, step := range steps {
			su := types.StepUpdate{StepName: step.Name, Status: recipe.StatusRunning}
			internalUpdates <- su
			updates <- su

			for _, spec := range step.Processes {
				proc, err := mgr.StartProcess(ctx, spec)
				if err != nil {
					fail := types.StepUpdate{StepName: step.Name, Status: recipe.StatusFailed, Err: err}
					internalUpdates <- fail
					updates <- fail
					return
				}
				// Send a message to the program to add the process tab
				r.program.Send(ui.AddProcessTabMsg{Name: spec.Name, Proc: proc})
			}

			done := types.StepUpdate{StepName: step.Name, Status: recipe.StatusDone}
			internalUpdates <- done
			updates <- done
		}
	}()

	if _, err := r.program.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}

func IsTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
