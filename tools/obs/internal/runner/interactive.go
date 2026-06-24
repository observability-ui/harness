package runner

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"obs/internal/process"
	"obs/internal/recipe"
	"obs/internal/ui"
)

func RunInteractive(ctx context.Context, mgr *process.Manager, prepare func() ([]*recipe.Step, error), updates chan<- recipe.StepUpdate) error {
	internalUpdates := make(chan recipe.StepUpdate, 100)
	retryCh := make(chan struct{}, 1)
	model := ui.NewModel(mgr, internalUpdates, retryCh)

	p := tea.NewProgram(model, tea.WithAltScreen())

	go func() {
		for {
			p.Send(ui.RequirementsCheckingMsg{})

			steps, err := prepare()
			if err != nil {
				p.Send(ui.RequirementsFailedMsg{Err: err})
				select {
				case <-retryCh:
					mgr.StopAll()
					continue
				case <-ctx.Done():
					return
				}
			}

			p.Send(ui.RequirementsPassedMsg{Steps: steps})

			cb := StepCallbacks{
				OnUpdate: func(u recipe.StepUpdate) { internalUpdates <- u },
				OnProcess: func(step *recipe.Step, spec recipe.ProcessSpec, proc *process.Process) {
					p.Send(ui.AddProcessTabMsg{StepName: step.Name, Name: spec.Name, Proc: proc})
				},
			}

			ExecuteSteps(ctx, mgr, steps, cb)

			select {
			case <-retryCh:
				mgr.StopAll()
				continue
			case <-ctx.Done():
				return
			}
		}
	}()

	if _, err := p.Run(); err != nil {
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
