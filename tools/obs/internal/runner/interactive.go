package runner

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"obs/internal/process"
	"obs/internal/recipe"
	"obs/internal/ui"
)

func RunInteractive(ctx context.Context, mgr *process.Manager, prepare func() ([]*recipe.Step, error), updates chan<- recipe.StepUpdate) error {
	internalUpdates := make(chan recipe.StepUpdate, 100)
	retryCh := make(chan struct{}, 1)
	model := ui.NewModel(mgr, nil, internalUpdates, retryCh)

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

			// Mark steps with dependencies as waiting
			for _, step := range steps {
				if len(step.DependsOn) > 0 {
					internalUpdates <- recipe.StepUpdate{StepName: step.Name, Status: recipe.StatusWaiting}
				}
			}

			ready := make(map[string]chan struct{})
			for _, step := range steps {
				ready[step.Name] = make(chan struct{})
			}

			for _, step := range steps {
				for _, dep := range step.DependsOn {
					if ch, ok := ready[dep]; ok {
						select {
						case <-ch:
						case <-ctx.Done():
							return
						}
					}
				}

				internalUpdates <- recipe.StepUpdate{StepName: step.Name, Status: recipe.StatusRunning}

				for _, spec := range step.Processes {
					proc, err := mgr.StartProcess(ctx, spec)
					if err != nil {
						internalUpdates <- recipe.StepUpdate{StepName: step.Name, Status: recipe.StatusFailed, Err: err}
						return
					}
					p.Send(ui.AddProcessTabMsg{StepName: step.Name, Name: spec.Name, Proc: proc})
				}

				internalUpdates <- recipe.StepUpdate{StepName: step.Name, Status: recipe.StatusStarted}

				if step.HasPorts() {
					var allPorts []int
					for _, spec := range step.Processes {
						allPorts = append(allPorts, spec.Ports...)
					}
					go func(name string, ports []int, ch chan struct{}) {
						for {
							if process.ProbePorts(ports) {
								internalUpdates <- recipe.StepUpdate{StepName: name, Status: recipe.StatusReady}
								close(ch)
								return
							}
							select {
							case <-time.After(time.Second):
							case <-ctx.Done():
								return
							}
						}
					}(step.Name, allPorts, ready[step.Name])
				} else {
					close(ready[step.Name])
				}
			}

			// Wait for retry signal or shutdown
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
