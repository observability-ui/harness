package strategy

import (
	"context"
	"strconv"

	"obs/internal/component"
	"obs/internal/runcontext"
)

type PodmanCompose struct{}

func (s *PodmanCompose) Name() string        { return "podman-compose" }
func (s *PodmanCompose) Requires() []string { return []string{"podman"} }

func (s *PodmanCompose) Run(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	composeFile := comp.Config["compose-file"]

	var ports []int
	for _, out := range comp.Outputs {
		if out.Name == "port" {
			if p, err := strconv.Atoi(out.Value); err == nil {
				ports = append(ports, p)
			}
		}
	}

	return &component.Step{
		Name:      comp.Name,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{{
			Name:    comp.Name,
			Command: "podman",
			Args:    []string{"compose", "-f", composeFile, "up"},
			Dir:     comp.Dir,
			Ports:   ports,
		}},
	}, nil
}
