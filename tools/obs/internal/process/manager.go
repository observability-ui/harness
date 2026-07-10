package process

import (
	"context"
	"fmt"
	"io"
	"sync"

	"obs/internal/task"
)

const defaultMaxLogLines = 10000

type Manager struct {
	mu        sync.RWMutex
	processes map[string]*Process
}

func NewManager() *Manager {
	return &Manager{
		processes: make(map[string]*Process),
	}
}

func (m *Manager) StartProcess(ctx context.Context, spec task.ProcessSpec, writers ...io.Writer) (*Process, error) {
	m.mu.Lock()
	if existing, ok := m.processes[spec.Name]; ok && existing.Running() {
		m.mu.Unlock()
		return nil, fmt.Errorf("process %q is already running", spec.Name)
	}

	proc := newProcess(spec, defaultMaxLogLines)
	m.processes[spec.Name] = proc
	m.mu.Unlock()

	if err := proc.Start(ctx, writers...); err != nil {
		m.mu.Lock()
		if m.processes[spec.Name] == proc {
			delete(m.processes, spec.Name)
		}
		m.mu.Unlock()
		return nil, fmt.Errorf("starting %q: %w", spec.Name, err)
	}
	return proc, nil
}

func (m *Manager) StopAll() {
	m.mu.RLock()
	procs := make([]*Process, 0, len(m.processes))
	for _, p := range m.processes {
		procs = append(procs, p)
	}
	m.mu.RUnlock()

	var wg sync.WaitGroup
	for _, p := range procs {
		wg.Add(1)
		go func(proc *Process) {
			defer wg.Done()
			proc.Stop()
		}(p)
	}
	wg.Wait()
}

func (m *Manager) RestartProcess(ctx context.Context, name string, writers ...io.Writer) (*Process, error) {
	m.mu.Lock()
	old, ok := m.processes[name]
	if !ok {
		m.mu.Unlock()
		return nil, fmt.Errorf("process %q not found", name)
	}
	spec := old.Spec
	m.mu.Unlock()

	old.Stop()

	m.mu.Lock()
	if current := m.processes[name]; current != old {
		m.mu.Unlock()
		return current, nil
	}
	proc := newProcess(spec, defaultMaxLogLines)
	proc.Output.Write([]byte("── restarting ──\n"))
	m.processes[name] = proc
	m.mu.Unlock()

	if err := proc.Start(ctx, writers...); err != nil {
		m.mu.Lock()
		if m.processes[name] == proc {
			delete(m.processes, name)
		}
		m.mu.Unlock()
		return nil, fmt.Errorf("restarting %q: %w", name, err)
	}
	return proc, nil
}

func (m *Manager) All() []*Process {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*Process, 0, len(m.processes))
	for _, p := range m.processes {
		result = append(result, p)
	}
	return result
}
