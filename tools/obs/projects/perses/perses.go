package perses

import (
	"context"

	"obs/internal/runcontext"
	"obs/internal/strategy"
	"obs/internal/task"
)

var BuildAPI = &task.Task{
	Name:     "perses-build",
	Dir:      "projects/perses",
	Strategy: strategy.MakeTarget("build-api"),
}

var Server = &task.Task{
	Name:      "perses-api",
	DependsOn: []string{"perses-build"},
	Dir:       "projects/perses",
	Ports:     []int{8080},
	Labels: map[string]string{
		"console-proxy-path": "/api/proxy/plugin/monitoring-console-plugin/perses/",
		"console-proxy-port": "8080",
	},
	Strategy: &PersesRunStrategy{},
}

func init() {
	task.Register(BuildAPI)
	task.Register(Server)
}

type PersesRunStrategy struct{}

func (s *PersesRunStrategy) Requires() []string { return []string{"go"} }

func (s *PersesRunStrategy) Execute(_ context.Context, t *task.Task, _ *runcontext.RunContext) (*task.Step, error) {
	return &task.Step{
		Name:      t.Name,
		Lifecycle: task.LifecycleLongRunning,
		DependsOn: t.DependsOn,
		Processes: []task.ProcessSpec{{
			Name:    "perses-api",
			Command: "bash",
			Args:    []string{"scripts/api_backend_dev.sh", "--e2e"},
			Dir:     t.Dir,
			Ports:   []int{8080},
		}},
	}, nil
}
