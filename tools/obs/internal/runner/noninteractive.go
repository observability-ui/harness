package runner

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"obs/internal/component"
	"obs/internal/process"
	"obs/internal/ui"

	"github.com/charmbracelet/lipgloss"
)

var prefixColors = []lipgloss.Color{
	lipgloss.Color("6"), // cyan
	lipgloss.Color("3"), // yellow
	lipgloss.Color("2"), // green
	lipgloss.Color("5"), // magenta
	lipgloss.Color("4"), // blue
	lipgloss.Color("1"), // red
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

func (r *NonInteractiveRunner) Run(ctx context.Context, mgr *process.Manager, steps []*component.Step, updates chan<- component.StepUpdate) error {
	colorIdx := 0

	cb := StepCallbacks{
		OnUpdate: func(u component.StepUpdate) { updates <- u },
		Writers: func(specName string) []io.Writer {
			pw := NewPrefixWriter(r.Out, specName, colorIdx)
			colorIdx++
			return []io.Writer{pw}
		},
	}

	launched, err := ExecuteSteps(ctx, mgr, steps, cb)
	if err != nil {
		close(updates)
		return err
	}

	type procResult struct {
		stepName string
		err      error
	}
	var wg sync.WaitGroup
	results := make(chan procResult, len(launched))

	for _, sp := range launched {
		wg.Add(1)
		go func(s StartedProc) {
			defer wg.Done()
			<-s.Proc.Wait()

			status := ui.MapProcessStatus(s.Proc.Status)
			updates <- component.StepUpdate{StepName: s.StepName, Status: status, Err: s.Proc.Err}
			var err error
			if status == component.StatusFailed {
				err = fmt.Errorf("process %q failed: %v", s.Proc.Spec.Name, s.Proc.Err)
			}
			results <- procResult{s.StepName, err}
		}(sp)
	}

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
