package strategy

import (
	"context"

	"obs/internal/runcontext"
	"obs/internal/task"
)

type npm struct {
	Cmd  string
	Args []string
}

func (s *npm) Requires() []string { return []string{"node", "npm"} }

func (s *npm) Execute(_ context.Context, t *task.Task, _ *runcontext.RunContext) (*task.Step, error) {
	return &task.Step{
		Name:      t.Name,
		Lifecycle: task.LifecycleOneShot,
		DependsOn: t.DependsOn,
		Processes: []task.ProcessSpec{{
			Name:    t.Name,
			Command: "npm",
			Args:    append([]string{s.Cmd}, s.Args...),
			Dir:     t.Dir,
		}},
	}, nil
}

func NPMRun(cmd string, args ...string) task.Strategy {
	return &npm{Cmd: cmd, Args: args}
}
