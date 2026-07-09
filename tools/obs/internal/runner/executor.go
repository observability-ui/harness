package runner

import (
	"context"
	"fmt"
	"io"
	"sync"

	"obs/internal/component"
	"obs/internal/process"
)

type stepErrors struct {
	mu   sync.Mutex
	errs map[string]error
}

func (se *stepErrors) set(name string, err error) {
	se.mu.Lock()
	defer se.mu.Unlock()
	se.errs[name] = err
}

func (se *stepErrors) get(name string) error {
	se.mu.Lock()
	defer se.mu.Unlock()
	return se.errs[name]
}

type StartedProc struct {
	StepName string
	Proc     *process.Process
}

type StepCallbacks struct {
	OnUpdate  func(component.StepUpdate)
	OnProcess func(step *component.Step, spec component.ProcessSpec, proc *process.Process)
	Writers   func(specName string) []io.Writer
}

func ExecuteSteps(ctx context.Context, mgr *process.Manager, steps []*component.Step, cb StepCallbacks) ([]StartedProc, error) {
	for _, step := range steps {
		if len(step.DependsOn) > 0 {
			cb.OnUpdate(component.StepUpdate{StepName: step.Name, Status: component.StatusWaiting})
		}
	}

	ready := make(map[string]chan struct{})
	stepErr := &stepErrors{errs: make(map[string]error)}
	for _, step := range steps {
		ready[step.Name] = make(chan struct{})
	}

	var launched []StartedProc

	for _, step := range steps {
		skip, err := waitDeps(ctx, step, ready, stepErr)
		if err != nil {
			return launched, err
		}
		if skip {
			cb.OnUpdate(component.StepUpdate{StepName: step.Name, Status: component.StatusSkipped})
			stepErr.set(step.Name, fmt.Errorf("dependency failed"))
			close(ready[step.Name])
			continue
		}

		cb.OnUpdate(component.StepUpdate{StepName: step.Name, Status: component.StatusRunning})

		var stepProcs []*process.Process
		for _, spec := range step.Processes {
			resolved, cleanup, err := ResolveSpec(spec)
			if err != nil {
				cb.OnUpdate(component.StepUpdate{StepName: step.Name, Status: component.StatusFailed, Err: err})
				return launched, fmt.Errorf("failed to resolve %q: %w", spec.Name, err)
			}

			var writers []io.Writer
			if cb.Writers != nil {
				writers = cb.Writers(spec.Name)
			}

			proc, err := mgr.StartProcess(ctx, resolved, writers...)
			if err != nil {
				if cleanup != nil {
					cleanup()
				}
				cb.OnUpdate(component.StepUpdate{StepName: step.Name, Status: component.StatusFailed, Err: err})
				return launched, fmt.Errorf("failed to start %q: %w", spec.Name, err)
			}
			if cleanup != nil {
				go func() { <-proc.Wait(); cleanup() }()
			}

			if cb.OnProcess != nil {
				cb.OnProcess(step, spec, proc)
			}

			launched = append(launched, StartedProc{StepName: step.Name, Proc: proc})
			stepProcs = append(stepProcs, proc)
		}

		cb.OnUpdate(component.StepUpdate{StepName: step.Name, Status: component.StatusStarted})

		go watchStepReady(ctx, step, stepProcs, ready[step.Name], stepErr, func(name string, status component.Status) {
			cb.OnUpdate(component.StepUpdate{StepName: name, Status: status})
		})
	}

	return launched, nil
}
