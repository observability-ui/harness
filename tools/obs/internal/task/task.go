package task

import (
	"context"
	"maps"

	"obs/internal/runcontext"
)

type RequiredFlag struct {
	Name  string
	Usage string
}

type Task struct {
	Name          string
	DependsOn     []string
	Dir           string
	Ports         []int
	RequiredFlags []RequiredFlag
	Labels        map[string]string
	Strategy      Strategy
}

type Strategy interface {
	Requires() []string
	Execute(ctx context.Context, t *Task, rc *runcontext.RunContext) (*Step, error)
}

var registry = make(map[string]*Task)

func Register(t *Task) {
	if _, exists := registry[t.Name]; exists {
		panic("duplicate task registration: " + t.Name)
	}
	registry[t.Name] = t
}

func Get(name string) (*Task, bool) {
	t, ok := registry[name]
	return t, ok
}

func All() map[string]*Task {
	cp := make(map[string]*Task, len(registry))
	maps.Copy(cp, registry)
	return cp
}

func ResetRegistry() {
	registry = make(map[string]*Task)
}
