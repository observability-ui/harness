package runcontext

import (
	"sync"
)

type RunContext struct {
	mu      sync.RWMutex
	outputs map[string]map[string]string
	tasks   map[string]bool
}

func New() *RunContext {
	return &RunContext{
		outputs: make(map[string]map[string]string),
		tasks:   make(map[string]bool),
	}
}

func (rc *RunContext) Set(taskName, key, value string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if rc.outputs[taskName] == nil {
		rc.outputs[taskName] = make(map[string]string)
	}
	rc.outputs[taskName][key] = value
}

func (rc *RunContext) Get(taskName, key string) string {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if m, ok := rc.outputs[taskName]; ok {
		return m[key]
	}
	return ""
}

func (rc *RunContext) SetTasks(names []string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.tasks = make(map[string]bool, len(names))
	for _, n := range names {
		rc.tasks[n] = true
	}
}

func (rc *RunContext) HasTask(name string) bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.tasks[name]
}
