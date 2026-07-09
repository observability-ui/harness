package strategy

import (
	"context"

	"obs/internal/component"
	"obs/internal/runcontext"
)

type Strategy interface {
	Requires() []string
	Execute(ctx context.Context, comp *component.Component, rc *runcontext.RunContext) (*component.Step, error)
}

var registry = map[string][]Strategy{}

func Register(componentName string, strategies ...Strategy) {
	registry[componentName] = append(registry[componentName], strategies...)
}

func Resolve(comp *component.Component) []Strategy {
	if s, ok := registry[comp.Name]; ok {
		return s
	}
	return resolveByConfig(comp)
}

func ResetRegistry() {
	registry = map[string][]Strategy{}
}
