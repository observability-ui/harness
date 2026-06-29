package process

import (
	"context"
	"fmt"
	"io"
	"sync"

	"obs/internal/recipe"
)

const DefaultMaxLogLines = 10000

type Manager struct {
	mu        sync.RWMutex
	processes map[string]*Process
}

func NewManager() *Manager {
	return &Manager{
		processes: make(map[string]*Process),
	}
}

func (m *Manager) StartProcess(ctx context.Context, spec recipe.ProcessSpec, writers ...io.Writer) (*Process, error) {
	m.mu.Lock()
	if existing, ok := m.processes[spec.Name]; ok && existing.Running() {
		m.mu.Unlock()
		return nil, fmt.Errorf("process %q is already running", spec.Name)
	}

	proc := NewProcess(spec, DefaultMaxLogLines)
	m.processes[spec.Name] = proc
	m.mu.Unlock()

	if err := proc.Start(ctx, writers...); err != nil {
		return nil, fmt.Errorf("starting %q: %w", spec.Name, err)
	}
	return proc, nil
}

func (m *Manager) Get(name string) (*Process, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.processes[name]
	return p, ok
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

func (m *Manager) StopProcess(name string) error {
	m.mu.RLock()
	p, ok := m.processes[name]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("process %q not found", name)
	}
	return p.Stop()
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
	proc := NewProcess(spec, DefaultMaxLogLines)
	proc.Output.Write([]byte("── restarting ──\n"))
	m.processes[name] = proc
	m.mu.Unlock()

	if err := proc.Start(ctx, writers...); err != nil {
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
