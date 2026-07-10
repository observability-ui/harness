package process

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"obs/internal/task"
)

type processStatus int

const (
	processPending processStatus = iota
	processRunning
	processDone
	processStopped
	processFailed
)

const shutdownTimeout = 5 * time.Second

type Process struct {
	Spec   task.ProcessSpec
	Status processStatus
	Err    error
	Output *ringBuffer

	cmd      *exec.Cmd
	mu       sync.Mutex
	done     chan struct{}
	doneOnce sync.Once
	stopping bool
}

func newProcess(spec task.ProcessSpec, maxLogLines int) *Process {
	return &Process{
		Spec:   spec,
		Status: processPending,
		Output: newRingBuffer(maxLogLines),
		done:   make(chan struct{}),
	}
}

func (p *Process) Start(ctx context.Context, extraWriters ...io.Writer) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Status != processPending {
		return fmt.Errorf("process %q already started", p.Spec.Name)
	}

	p.cmd = exec.Command(p.Spec.Command, p.Spec.Args...)
	if p.Spec.Dir != "" {
		p.cmd.Dir = p.Spec.Dir
	}
	if len(p.Spec.Env) > 0 {
		overrides := make(map[string]bool, len(p.Spec.Env))
		for k := range p.Spec.Env {
			overrides[k] = true
		}
		for _, e := range os.Environ() {
			if k, _, ok := strings.Cut(e, "="); ok && overrides[k] {
				continue
			}
			p.cmd.Env = append(p.cmd.Env, e)
		}
		for k, v := range p.Spec.Env {
			p.cmd.Env = append(p.cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	// Put child in its own process group for clean shutdown
	p.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	writers := []io.Writer{p.Output}
	writers = append(writers, extraWriters...)
	mw := io.MultiWriter(writers...)
	p.cmd.Stdout = mw
	p.cmd.Stderr = mw

	if err := p.cmd.Start(); err != nil {
		p.Status = processFailed
		p.Err = err
		p.doneOnce.Do(func() { close(p.done) })
		return err
	}

	p.Status = processRunning

	go func() {
		// Stop the process group when context is cancelled
		select {
		case <-ctx.Done():
			p.Stop()
		case <-p.done:
		}
	}()

	go func() {
		err := p.cmd.Wait()
		p.mu.Lock()
		if p.stopping {
			p.Status = processStopped
		} else if err != nil {
			p.Status = processFailed
			p.Err = err
		} else {
			p.Status = processDone
		}
		p.mu.Unlock()
		p.doneOnce.Do(func() { close(p.done) })
	}()

	return nil
}

func (p *Process) Stop() error {
	p.mu.Lock()
	if p.cmd == nil || p.cmd.Process == nil || p.Status != processRunning {
		p.mu.Unlock()
		return nil
	}
	p.stopping = true
	pid := p.cmd.Process.Pid
	p.mu.Unlock()

	// SIGINT to the process group
	syscall.Kill(-pid, syscall.SIGINT)

	select {
	case <-p.done:
		return nil
	case <-time.After(shutdownTimeout):
		syscall.Kill(-pid, syscall.SIGKILL)
		select {
		case <-p.done:
			return nil
		case <-time.After(3 * time.Second):
			p.mu.Lock()
			if p.Status == processRunning {
				p.Status = processStopped
			}
			p.mu.Unlock()
			p.doneOnce.Do(func() { close(p.done) })
			return fmt.Errorf("process %d did not exit after SIGKILL", pid)
		}
	}
}

func (p *Process) Wait() <-chan struct{} {
	return p.done
}

func (p *Process) Running() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Status == processRunning
}

func (p *Process) GetStatus() processStatus {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Status
}

func (p *Process) GetErr() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Err
}

func (p *Process) PID() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.cmd != nil && p.cmd.Process != nil {
		return p.cmd.Process.Pid
	}
	return 0
}
