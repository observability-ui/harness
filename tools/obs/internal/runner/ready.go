package runner

import (
	"context"
	"fmt"
	"time"

	"obs/internal/task"
	"obs/internal/process"
)

func waitDeps(ctx context.Context, step *task.Step, ready map[string]chan struct{}, stepErr *stepErrors) (skip bool, err error) {
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

func watchStepReady(ctx context.Context, step *task.Step, procs []*process.Process, ch chan struct{}, stepErr *stepErrors, onReady func(string, task.Status)) {
	defer close(ch)

	allDone := process.WaitAll(procs)

	if step.Lifecycle == task.LifecycleLongRunning {
		var allPorts []int
		for _, spec := range step.Processes {
			allPorts = append(allPorts, spec.Ports...)
		}

		for {
			if process.ProbePorts(allPorts) {
				onReady(step.Name, task.StatusReady)
				break
			}
			select {
			case <-time.After(time.Second):
			case <-allDone:
				stepErr.set(step.Name, fmt.Errorf("process exited before ports were ready"))
				onReady(step.Name, task.StatusFailed)
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

func reportExitStatus(stepName string, procs []*process.Process, stepErr *stepErrors, onReady func(string, task.Status)) {
	status := aggregateExitStatus(procs)
	if status == task.StatusFailed || status == task.StatusStopped {
		stepErr.set(stepName, fmt.Errorf("step %q: %s", stepName, status))
	}
	onReady(stepName, status)
}

func aggregateExitStatus(procs []*process.Process) task.Status {
	worst := task.StatusDone
	for _, proc := range procs {
		s := process.MapExitStatus(proc)
		if s == task.StatusFailed {
			return task.StatusFailed
		}
		if s == task.StatusStopped {
			worst = task.StatusStopped
		}
	}
	return worst
}
