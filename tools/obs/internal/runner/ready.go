package runner

import (
	"context"
	"fmt"
	"time"

	"obs/internal/process"
	"obs/internal/recipe"
)

func waitDeps(ctx context.Context, step *recipe.Step, ready map[string]chan struct{}, stepErr map[string]error) (skip bool, err error) {
	for _, dep := range step.DependsOn {
		if ch, ok := ready[dep]; ok {
			select {
			case <-ch:
				if stepErr[dep] != nil {
					return true, nil
				}
			case <-ctx.Done():
				return false, ctx.Err()
			}
		}
	}
	return false, nil
}

func watchStepReady(ctx context.Context, step *recipe.Step, procs []*process.Process, ch chan struct{}, stepErr map[string]error, onReady func(string, recipe.Status)) {
	if step.HasPorts() {
		var allPorts []int
		for _, spec := range step.Processes {
			allPorts = append(allPorts, spec.Ports...)
		}

		allDone := make(chan struct{})
		go func() {
			for _, proc := range procs {
				<-proc.Wait()
			}
			close(allDone)
		}()

		for {
			if process.ProbePorts(allPorts) {
				onReady(step.Name, recipe.StatusReady)
				close(ch)
				return
			}
			select {
			case <-time.After(time.Second):
			case <-allDone:
				stepErr[step.Name] = fmt.Errorf("process exited before ports were ready")
				close(ch)
				return
			case <-ctx.Done():
				return
			}
		}
	}

	for _, proc := range procs {
		select {
		case <-proc.Wait():
		case <-ctx.Done():
			return
		}
	}
	for _, proc := range procs {
		if proc.Status == process.ProcessFailed {
			stepErr[step.Name] = fmt.Errorf("process %q failed: %v", proc.Spec.Name, proc.Err)
			close(ch)
			return
		}
	}
	close(ch)
}
