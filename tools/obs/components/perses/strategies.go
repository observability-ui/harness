package perses

import (
	"context"

	"obs/internal/component"
	"obs/internal/runcontext"
	"obs/internal/strategy"
)

type PersesRunStrategy struct{}

func (s *PersesRunStrategy) Name() string        { return "perses-local" }
func (s *PersesRunStrategy) Requires() []string { return []string{"go"} }

func (s *PersesRunStrategy) Run(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	return &component.Step{
		Name:      comp.Name,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{{
			Name:    "perses-api",
			Command: "./bin/perses",
			Args:    []string{"--config", "{{path:local-config}}"},
			Dir:     comp.Dir,
			Ports:   []int{8080},
			Files: map[string]component.FileRef{
				"local-config": {FS: filesFS, Path: "files/local-config.yaml"},
			},
		}},
	}, nil
}

func init() {
	strategy.RegisterSelector(func(comp *component.Component, mode string) (strategy.BuildStrategy, strategy.RunStrategy) {
		if comp.Name == Server.Name {
			return nil, &PersesRunStrategy{}
		}
		return nil, nil
	})
}
