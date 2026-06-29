package runner

import (
	"context"
	"fmt"
	"io"

	"obs/internal/process"
	"obs/internal/recipe"
)

type StartedProc struct {
	StepName string
	Proc     *process.Process
}

type StepCallbacks struct {
	OnUpdate  func(recipe.StepUpdate)
	OnProcess func(step *recipe.Step, spec recipe.ProcessSpec, proc *process.Process)
	Writers   func(specName string) []io.Writer
}

func ExecuteSteps(ctx context.Context, mgr *process.Manager, steps []*recipe.Step, cb StepCallbacks) ([]StartedProc, error) {
	for _, step := range steps {
		if len(step.DependsOn) > 0 {
			cb.OnUpdate(recipe.StepUpdate{StepName: step.Name, Status: recipe.StatusWaiting})
		}
	}

	ready := make(map[string]chan struct{})
	stepErr := make(map[string]error)
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
			cb.OnUpdate(recipe.StepUpdate{StepName: step.Name, Status: recipe.StatusSkipped})
			stepErr[step.Name] = fmt.Errorf("dependency failed")
			close(ready[step.Name])
			continue
		}

		cb.OnUpdate(recipe.StepUpdate{StepName: step.Name, Status: recipe.StatusRunning})

		var stepProcs []*process.Process
		for _, spec := range step.Processes {
			resolved, cleanup, err := recipe.ResolveSpec(spec)
			if err != nil {
				cb.OnUpdate(recipe.StepUpdate{StepName: step.Name, Status: recipe.StatusFailed, Err: err})
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
				cb.OnUpdate(recipe.StepUpdate{StepName: step.Name, Status: recipe.StatusFailed, Err: err})
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

		cb.OnUpdate(recipe.StepUpdate{StepName: step.Name, Status: recipe.StatusStarted})

		go watchStepReady(ctx, step, stepProcs, ready[step.Name], stepErr, func(name string, status recipe.Status) {
			cb.OnUpdate(recipe.StepUpdate{StepName: name, Status: status})
		})
	}

	return launched, nil
}
