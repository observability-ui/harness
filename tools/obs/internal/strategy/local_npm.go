package strategy

import (
	"context"
	"path/filepath"

	"obs/internal/component"
	"obs/internal/runcontext"
)

type LocalNPMInstall struct{}

func (s *LocalNPMInstall) Name() string        { return "local-npm-install" }
func (s *LocalNPMInstall) Requires() []string { return []string{"node", "npm"} }

func (s *LocalNPMInstall) Build(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	dir := comp.Dir
	if sub := comp.Config["subdir"]; sub != "" {
		dir = filepath.Join(dir, sub)
	}
	cmd := comp.Config["npm-cmd"]
	if cmd == "" {
		cmd = "install"
	}

	return &component.Step{
		Name:      comp.Name,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{{
			Name:    comp.Name,
			Command: "npm",
			Args:    []string{cmd},
			Dir:     dir,
		}},
	}, nil
}
