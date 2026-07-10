package process

import (
	"time"

	"obs/internal/task"
)

func WaitAll(procs []*Process) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		for _, p := range procs {
			<-p.Wait()
		}
		close(ch)
	}()
	return ch
}

func WaitForReady(proc *Process, ports []int) task.Status {
	if len(ports) == 0 {
		<-proc.Wait()
		return MapExitStatus(proc)
	}
	for {
		if ProbePorts(ports) {
			return task.StatusReady
		}
		select {
		case <-time.After(time.Second):
		case <-proc.Wait():
			return task.StatusFailed
		}
	}
}

func WaitForExit(proc *Process) task.Status {
	<-proc.Wait()
	return MapExitStatus(proc)
}

func MapExitStatus(proc *Process) task.Status {
	switch proc.GetStatus() {
	case processFailed:
		return task.StatusFailed
	case processStopped:
		return task.StatusStopped
	default:
		return task.StatusDone
	}
}
