package process

import (
	"time"

	"obs/internal/component"
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

func WaitForReady(proc *Process, ports []int) component.Status {
	if len(ports) == 0 {
		<-proc.Wait()
		return MapExitStatus(proc)
	}
	for {
		if ProbePorts(ports) {
			return component.StatusReady
		}
		select {
		case <-time.After(time.Second):
		case <-proc.Wait():
			return component.StatusFailed
		}
	}
}

func WaitForExit(proc *Process) component.Status {
	<-proc.Wait()
	return MapExitStatus(proc)
}

func MapExitStatus(proc *Process) component.Status {
	switch proc.GetStatus() {
	case ProcessFailed:
		return component.StatusFailed
	case ProcessStopped:
		return component.StatusStopped
	default:
		return component.StatusDone
	}
}
