package strategy

import (
	"context"
	"fmt"

	"obs/internal/component"
	"obs/internal/runcontext"
)

type LocalMakeRun struct{}

func (s *LocalMakeRun) Requires() []string { return []string{"make"} }

func (s *LocalMakeRun) Execute(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	target := comp.Config["make-target"]
	if target == "" {
		return nil, fmt.Errorf("component %q: missing make-target config", comp.Name)
	}

	ports := comp.Ports
	lifecycle := component.LifecycleOneShot
	if len(ports) > 0 {
		lifecycle = component.LifecycleLongRunning
	}

	return &component.Step{
		Name:      comp.Name,
		Lifecycle: lifecycle,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{{
			Name:    comp.Name,
			Command: "make",
			Args:    []string{target},
			Dir:     comp.Dir,
			Ports:   ports,
		}},
	}, nil
}
