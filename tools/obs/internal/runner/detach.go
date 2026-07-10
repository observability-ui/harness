package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"obs/internal/task"
	"obs/internal/process"
	"obs/internal/state"
)

type DetachRunner struct {
	store *state.Store
	lock  *state.Lock
}

func NewDetach(stateDir string) *DetachRunner {
	return &DetachRunner{
		store: state.NewStore(stateDir),
		lock:  state.NewLock(stateDir),
	}
}

func (r *DetachRunner) Run(ctx context.Context, mgr *process.Manager, steps []*task.Step, updates chan<- task.StepUpdate) error {
	if err := r.lock.Acquire(); err != nil {
		return err
	}

	inner := NewNonInteractive(io.Discard)
	errCh := make(chan error, 1)
	go func() {
		err := inner.Run(ctx, mgr, steps, updates)
		errCh <- err
		r.lock.Release()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("background run failed: %w", err)
		}
	case <-time.After(1 * time.Second):
	}

	var procs []state.ProcessState
	for _, p := range mgr.All() {
		procs = append(procs, state.ProcessState{
			Name: p.Spec.Name,
			PID:  p.PID(),
		})
		if err := r.store.WritePID(p.Spec.Name, p.PID()); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to write PID for %s: %v\n", p.Spec.Name, err)
		}
	}
	if err := r.store.Save(&state.RunState{Processes: procs}); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to save state: %v\n", err)
	}

	fmt.Println("Processes started in background:")
	state.PrintStatus(&state.RunState{Processes: procs})
	return nil
}
