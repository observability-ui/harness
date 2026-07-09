package strategy

import (
	"context"
	"fmt"

	"obs/internal/component"
	"obs/internal/runcontext"
)

type PodmanCompose struct{}

func (s *PodmanCompose) Requires() []string { return []string{"podman"} }

func (s *PodmanCompose) Execute(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	composeFile := comp.Config["compose-file"]
	if composeFile == "" {
		return nil, fmt.Errorf("component %q: missing compose-file config", comp.Name)
	}

	return &component.Step{
		Name:      comp.Name,
		Lifecycle: component.LifecycleLongRunning,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{{
			Name:    comp.Name,
			Command: "podman",
			Args:    []string{"compose", "-f", composeFile, "up"},
			Dir:     comp.Dir,
			Ports:   comp.Ports,
		}},
	}, nil
}
