package strategy

import (
	"context"
	"fmt"

	"obs/internal/runcontext"
	"obs/internal/task"
)

type podmanCompose struct {
	File string
}

func (s *podmanCompose) Requires() []string { return []string{"podman"} }

func (s *podmanCompose) Execute(_ context.Context, t *task.Task, _ *runcontext.RunContext) (*task.Step, error) {
	if s.File == "" {
		return nil, fmt.Errorf("task %q: missing compose file", t.Name)
	}

	return &task.Step{
		Name:      t.Name,
		Lifecycle: task.LifecycleLongRunning,
		DependsOn: t.DependsOn,
		Processes: []task.ProcessSpec{{
			Name:    t.Name,
			Command: "podman",
			Args:    []string{"compose", "-f", s.File, "up"},
			Dir:     t.Dir,
			Ports:   t.Ports,
		}},
	}, nil
}

func Compose(file string) task.Strategy {
	return &podmanCompose{File: file}
}
