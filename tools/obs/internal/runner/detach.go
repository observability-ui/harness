package runner

import (
	"context"
	"fmt"
	"io"
	"time"

	"obs/internal/process"
	"obs/internal/recipe"
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

func (r *DetachRunner) Run(ctx context.Context, mgr *process.Manager, steps []*recipe.Step, updates chan<- recipe.StepUpdate) error {
	if err := r.lock.Acquire(); err != nil {
		return err
	}

	inner := NewNonInteractive(io.Discard)
	go inner.Run(ctx, mgr, steps, updates)

	time.Sleep(500 * time.Millisecond)

	var procs []state.ProcessState
	for _, p := range mgr.All() {
		procs = append(procs, state.ProcessState{
			Name: p.Spec.Name,
			PID:  p.PID(),
		})
		r.store.WritePID(p.Spec.Name, p.PID())
	}
	r.store.Save(&state.RunState{Processes: procs})

	fmt.Println("Processes started in background:")
	state.PrintStatus(&state.RunState{Processes: procs})
	return nil
}
