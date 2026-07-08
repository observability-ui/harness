package strategy

import (
	"context"
	"strconv"
	"strings"

	"obs/internal/component"
	"obs/internal/runcontext"
)

type LocalMakeRun struct{}

func (s *LocalMakeRun) Name() string        { return "local-make-run" }
func (s *LocalMakeRun) Requires() []string { return []string{"make"} }

func (s *LocalMakeRun) Run(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	target := comp.Config["make-target"]

	var ports []int
	for _, out := range comp.Outputs {
		if out.Name == "port" {
			if p, err := strconv.Atoi(out.Value); err == nil {
				ports = append(ports, p)
			}
		}
	}

	env := make(map[string]string)
	for k, v := range comp.Config {
		if after, ok := strings.CutPrefix(k, "env."); ok {
			env[after] = v
		}
	}

	return &component.Step{
		Name:      comp.Name,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{{
			Name:    comp.Name,
			Command: "make",
			Args:    []string{target},
			Dir:     comp.Dir,
			Ports:   ports,
			Env:     env,
		}},
	}, nil
}
