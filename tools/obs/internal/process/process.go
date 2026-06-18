package process

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"obs/internal/recipe"
)

type ProcessStatus int

const (
	ProcessPending ProcessStatus = iota
	ProcessRunning
	ProcessStopped
	ProcessFailed
)

const ShutdownTimeout = 5 * time.Second

type Process struct {
	Spec   recipe.ProcessSpec
	Status ProcessStatus
	Err    error
	Output *RingBuffer

	cmd  *exec.Cmd
	mu   sync.Mutex
	done chan struct{}
}

func NewProcess(spec recipe.ProcessSpec, maxLogLines int) *Process {
	return &Process{
		Spec:   spec,
		Status: ProcessPending,
		Output: NewRingBuffer(maxLogLines),
		done:   make(chan struct{}),
	}
}

func (p *Process) Start(ctx context.Context, extraWriters ...io.Writer) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.cmd = exec.CommandContext(ctx, p.Spec.Command, p.Spec.Args...)
	if p.Spec.Dir != "" {
		p.cmd.Dir = p.Spec.Dir
	}
	if len(p.Spec.Env) > 0 {
		p.cmd.Env = os.Environ()
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
		p.Status = ProcessFailed
		p.Err = err
		close(p.done)
		return err
	}

	p.Status = ProcessRunning

	go func() {
		err := p.cmd.Wait()
		p.mu.Lock()
		if err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) && exitErr.ExitCode() == -1 {
				p.Status = ProcessStopped
			} else {
				p.Status = ProcessFailed
				p.Err = err
			}
		} else {
			p.Status = ProcessStopped
		}
		p.mu.Unlock()
		close(p.done)
	}()

	return nil
}

func (p *Process) Stop() error {
	p.mu.Lock()
	if p.cmd == nil || p.cmd.Process == nil || p.Status != ProcessRunning {
		p.mu.Unlock()
		return nil
	}
	pid := p.cmd.Process.Pid
	p.mu.Unlock()

	// SIGINT to the process group
	syscall.Kill(-pid, syscall.SIGINT)

	select {
	case <-p.done:
		return nil
	case <-time.After(ShutdownTimeout):
		syscall.Kill(-pid, syscall.SIGKILL)
		<-p.done
		return nil
	}
}

func (p *Process) Wait() <-chan struct{} {
	return p.done
}

func (p *Process) Running() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Status == ProcessRunning
}

func (p *Process) PID() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.cmd != nil && p.cmd.Process != nil {
		return p.cmd.Process.Pid
	}
	return 0
}
