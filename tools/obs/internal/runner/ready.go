package runner

import (
	"context"
	"fmt"
	"time"

	"obs/internal/component"
	"obs/internal/process"
)

func waitDeps(ctx context.Context, step *component.Step, ready map[string]chan struct{}, stepErr *stepErrors) (skip bool, err error) {
	for _, dep := range step.DependsOn {
		if ch, ok := ready[dep]; ok {
			select {
			case <-ch:
				if stepErr.get(dep) != nil {
					return true, nil
				}
			case <-ctx.Done():
				return false, ctx.Err()
			}
		}
	}
	return false, nil
}

func watchStepReady(ctx context.Context, step *component.Step, procs []*process.Process, ch chan struct{}, stepErr *stepErrors, onReady func(string, component.Status)) {
	defer close(ch)

	allDone := process.WaitAll(procs)

	if step.Lifecycle == component.LifecycleLongRunning {
		var allPorts []int
		for _, spec := range step.Processes {
			allPorts = append(allPorts, spec.Ports...)
		}

		for {
			if process.ProbePorts(allPorts) {
				onReady(step.Name, component.StatusReady)
				break
			}
			select {
			case <-time.After(time.Second):
			case <-allDone:
				stepErr.set(step.Name, fmt.Errorf("process exited before ports were ready"))
				onReady(step.Name, component.StatusFailed)
				return
			case <-ctx.Done():
				<-allDone
				reportExitStatus(step.Name, procs, stepErr, onReady)
				return
			}
		}

		<-allDone
		reportExitStatus(step.Name, procs, stepErr, onReady)
		return
	}

	for _, proc := range procs {
		<-proc.Wait()
	}
	reportExitStatus(step.Name, procs, stepErr, onReady)
}

func reportExitStatus(stepName string, procs []*process.Process, stepErr *stepErrors, onReady func(string, component.Status)) {
	status := aggregateExitStatus(procs)
	if status == component.StatusFailed || status == component.StatusStopped {
		stepErr.set(stepName, fmt.Errorf("step %q: %s", stepName, status))
	}
	onReady(stepName, status)
}

func aggregateExitStatus(procs []*process.Process) component.Status {
	worst := component.StatusDone
	for _, proc := range procs {
		s := process.MapExitStatus(proc)
		if s == component.StatusFailed {
			return component.StatusFailed
		}
		if s == component.StatusStopped {
			worst = component.StatusStopped
		}
	}
	return worst
}
