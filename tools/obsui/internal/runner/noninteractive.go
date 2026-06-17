package runner

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"obsui/internal/process"
	"obsui/internal/recipe"
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

func (r *NonInteractiveRunner) Run(ctx context.Context, mgr *process.Manager, steps []*recipe.Step, updates chan<- StepUpdate) error {
	defer close(updates)

	colorIdx := 0
	var allProcs []*process.Process

	for _, step := range steps {
		updates <- StepUpdate{StepName: step.Name, Status: recipe.StatusRunning}

		for _, spec := range step.Processes {
			pw := NewPrefixWriter(r.Out, spec.Name, colorIdx)
			colorIdx++

			proc, err := mgr.StartProcess(ctx, spec, pw)
			if err != nil {
				updates <- StepUpdate{StepName: step.Name, Status: recipe.StatusFailed, Err: err}
				return fmt.Errorf("failed to start %q: %w", spec.Name, err)
			}
			allProcs = append(allProcs, proc)
		}

		updates <- StepUpdate{StepName: step.Name, Status: recipe.StatusDone}
	}

	// Wait for all processes to complete (for finite recipes like deploy)
	// For long-running recipes (start), this blocks until ctx is cancelled
	for _, proc := range allProcs {
		select {
		case <-proc.Wait():
			if proc.Status == process.ProcessFailed {
				return fmt.Errorf("process %q failed: %v", proc.Spec.Name, proc.Err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
