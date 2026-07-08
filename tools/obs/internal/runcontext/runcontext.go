package runcontext

import (
	"maps"
	"sync"
)

type RunContext struct {
	mu      sync.RWMutex
	outputs map[string]map[string]string
	mode    string
}

func New(mode string) *RunContext {
	return &RunContext{
		outputs: make(map[string]map[string]string),
		mode:    mode,
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

func (rc *RunContext) GetAll(key string) []map[string]string {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	var results []map[string]string
	for _, m := range rc.outputs {
		if _, ok := m[key]; ok {
			cp := make(map[string]string, len(m))
			maps.Copy(cp, m)
			results = append(results, cp)
		}
	}
	return results
}

func (rc *RunContext) Mode() string {
	return rc.mode
}
