package strategy

import (
	"context"
	"strings"

	"obs/internal/component"
	"obs/internal/runcontext"
)

type OCCommand struct{}

func (s *OCCommand) Name() string        { return "oc-command" }
func (s *OCCommand) Requires() []string { return []string{"oc"} }

func (s *OCCommand) Run(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	args := strings.Fields(comp.Config["oc-args"])

	return &component.Step{
		Name:      comp.Name,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{{
			Name:    comp.Name,
			Command: "oc",
			Args:    args,
		}},
	}, nil
}
