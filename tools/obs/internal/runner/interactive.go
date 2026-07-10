package runner

import (
	"context"
	"fmt"
	"os"

	"obs/internal/task"
	"obs/internal/process"
	"obs/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func RunInteractive(ctx context.Context, mgr *process.Manager, prepare func() ([]*task.Step, []task.ProjectInfo, error), inputValues map[string]string) error {
	innerCtx, innerCancel := context.WithCancel(ctx)
	defer innerCancel()

	internalUpdates := make(chan task.StepUpdate, 100)
	retryCh := make(chan struct{}, 1)
	model := ui.NewModel(innerCtx, mgr, internalUpdates, retryCh, inputValues)

	p := tea.NewProgram(model, tea.WithAltScreen())

	go func() {
		for {
			p.Send(ui.RequirementsCheckingMsg{})

			steps, projects, err := prepare()
			if err != nil {
				p.Send(ui.RequirementsFailedMsg{Err: err})
				select {
				case <-retryCh:
					mgr.StopAll()
					continue
				case <-innerCtx.Done():
					return
				}
			}

			p.Send(ui.RequirementsPassedMsg{Steps: steps, Projects: projects})

			cb := StepCallbacks{
				OnUpdate: func(u task.StepUpdate) { internalUpdates <- u },
				OnProcess: func(step *task.Step, spec task.ProcessSpec, proc *process.Process) {
					p.Send(ui.AddProcessTabMsg{StepName: step.Name, Name: spec.Name, Proc: proc})
				},
			}

			if _, err := ExecuteSteps(innerCtx, mgr, steps, cb); err != nil {
				p.Send(ui.RequirementsFailedMsg{Err: err})
			}

			select {
			case <-retryCh:
				mgr.StopAll()
				continue
			case <-innerCtx.Done():
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
