package runner

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"obs/internal/process"
	"obs/internal/recipe"
)

var prefixColors = []lipgloss.Color{
	lipgloss.Color("6"),  // cyan
	lipgloss.Color("3"),  // yellow
	lipgloss.Color("2"),  // green
	lipgloss.Color("5"),  // magenta
	lipgloss.Color("4"),  // blue
	lipgloss.Color("1"),  // red
}

type PrefixWriter struct {
	w       io.Writer
	prefix  string
	partial string
	mu      sync.Mutex
}

func NewPrefixWriter(w io.Writer, name string, colorIdx int) *PrefixWriter {
	color := prefixColors[colorIdx%len(prefixColors)]
	style := lipgloss.NewStyle().Foreground(color)
	prefix := style.Render(fmt.Sprintf("%-20s | ", name))
	return &PrefixWriter{w: w, prefix: prefix}
}

func (pw *PrefixWriter) Write(p []byte) (int, error) {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	text := pw.partial + string(p)
	lines := strings.Split(text, "\n")
	pw.partial = lines[len(lines)-1]

	for _, line := range lines[:len(lines)-1] {
		fmt.Fprintf(pw.w, "%s%s\n", pw.prefix, line)
	}
	return len(p), nil
}

type NonInteractiveRunner struct {
	Out io.Writer
}

func NewNonInteractive(out io.Writer) *NonInteractiveRunner {
	return &NonInteractiveRunner{Out: out}
}

func (r *NonInteractiveRunner) Run(ctx context.Context, mgr *process.Manager, steps []*recipe.Step, updates chan<- recipe.StepUpdate) error {
	type stepProc struct {
		stepName string
		proc     *process.Process
	}

	// Mark steps with dependencies as waiting
	for _, step := range steps {
		if len(step.DependsOn) > 0 {
			updates <- recipe.StepUpdate{StepName: step.Name, Status: recipe.StatusWaiting}
		}
	}

	colorIdx := 0
	var launched []stepProc
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
					close(updates)
					return ctx.Err()
				}
			}
		}

		updates <- recipe.StepUpdate{StepName: step.Name, Status: recipe.StatusRunning}

		for _, spec := range step.Processes {
			pw := NewPrefixWriter(r.Out, spec.Name, colorIdx)
			colorIdx++

			proc, err := mgr.StartProcess(ctx, spec, pw)
			if err != nil {
				updates <- recipe.StepUpdate{StepName: step.Name, Status: recipe.StatusFailed, Err: err}
				close(updates)
				return fmt.Errorf("failed to start %q: %w", spec.Name, err)
			}
			launched = append(launched, stepProc{stepName: step.Name, proc: proc})
		}

		updates <- recipe.StepUpdate{StepName: step.Name, Status: recipe.StatusStarted}

		if !step.HasPorts() {
			close(ready[step.Name])
		}
	}

	// Monitor with WaitGroup — channel stays open until all goroutines finish
	type procResult struct {
		stepName string
		err      error
	}
	var wg sync.WaitGroup
	results := make(chan procResult, len(launched))

	for _, sp := range launched {
		wg.Add(1)
		go func(s stepProc) {
			defer wg.Done()
			ports := s.proc.Spec.Ports

			if len(ports) > 0 {
				for {
					if process.ProbePorts(ports) {
						updates <- recipe.StepUpdate{StepName: s.stepName, Status: recipe.StatusReady}
						if ch, ok := ready[s.stepName]; ok {
							select {
							case <-ch:
							default:
								close(ch)
							}
						}
						break
					}
					select {
					case <-time.After(time.Second):
					case <-s.proc.Wait():
						goto exited
					}
				}
			}

			<-s.proc.Wait()
		exited:
			switch s.proc.Status {
			case process.ProcessFailed:
				updates <- recipe.StepUpdate{StepName: s.stepName, Status: recipe.StatusFailed, Err: s.proc.Err}
				results <- procResult{s.stepName, fmt.Errorf("process %q failed: %v", s.proc.Spec.Name, s.proc.Err)}
			case process.ProcessStopped:
				updates <- recipe.StepUpdate{StepName: s.stepName, Status: recipe.StatusStopped}
				results <- procResult{s.stepName, nil}
			default:
				updates <- recipe.StepUpdate{StepName: s.stepName, Status: recipe.StatusDone}
				if ch, ok := ready[s.stepName]; ok {
					select {
					case <-ch:
					default:
						close(ch)
					}
				}
				results <- procResult{s.stepName, nil}
			}
		}(sp)
	}

	// Wait for goroutines to finish so we can close the channel safely
	allDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(allDone)
	}()

	var firstErr error
	for remaining := len(launched); remaining > 0; {
		select {
		case r := <-results:
			remaining--
			if r.err != nil && firstErr == nil {
				firstErr = r.err
			}
		case <-ctx.Done():
			mgr.StopAll()
			<-allDone
			close(updates)
			return nil
		}
	}
	close(updates)
	return firstErr
}
