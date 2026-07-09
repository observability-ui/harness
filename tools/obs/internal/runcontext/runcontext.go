package runcontext

import (
	"sync"
)

type RunContext struct {
	mu         sync.RWMutex
	outputs    map[string]map[string]string
	components map[string]bool
}

func New() *RunContext {
	return &RunContext{
		outputs:    make(map[string]map[string]string),
		components: make(map[string]bool),
	}
}

func (rc *RunContext) Set(component, key, value string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if rc.outputs[component] == nil {
		rc.outputs[component] = make(map[string]string)
	}
	rc.outputs[component][key] = value
}

func (rc *RunContext) Get(component, key string) string {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if m, ok := rc.outputs[component]; ok {
		return m[key]
	}
	return ""
}

func (rc *RunContext) SetComponents(names []string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.components = make(map[string]bool, len(names))
	for _, n := range names {
		rc.components[n] = true
	}
}

func (rc *RunContext) HasComponent(name string) bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.components[name]
}
