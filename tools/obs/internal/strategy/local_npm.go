package strategy

import (
	"context"

	"obs/internal/component"
	"obs/internal/runcontext"
)

type LocalNPM struct {
	Cmd  string
	Args []string
}

func (s *LocalNPM) Requires() []string { return []string{"node", "npm"} }

func (s *LocalNPM) Execute(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	return &component.Step{
		Name:      comp.Name,
		Lifecycle: component.LifecycleOneShot,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{{
			Name:    comp.Name,
			Command: "npm",
			Args:    append([]string{s.Cmd}, s.Args...),
			Dir:     comp.Dir,
		}},
	}, nil
}
