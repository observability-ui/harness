package strategy

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"obs/internal/runcontext"
	"obs/internal/task"
)

type makeRun struct {
	Target string
}

func (s *makeRun) Requires() []string { return []string{"make"} }

func (s *makeRun) Execute(_ context.Context, t *task.Task, _ *runcontext.RunContext) (*task.Step, error) {
	if s.Target == "" {
		return nil, fmt.Errorf("task %q: missing make target", t.Name)
	}

	target, err := resolveTarget(t.Dir, s.Target)
	if err != nil {
		return nil, fmt.Errorf("task %q: %w", t.Name, err)
	}

	ports := t.Ports
	lifecycle := task.LifecycleOneShot
	if len(ports) > 0 {
		lifecycle = task.LifecycleLongRunning
	}

	return &task.Step{
		Name:      t.Name,
		Lifecycle: lifecycle,
		DependsOn: t.DependsOn,
		Processes: []task.ProcessSpec{{
			Name:    t.Name,
			Command: "make",
			Args:    []string{target},
			Dir:     t.Dir,
			Ports:   ports,
		}},
	}, nil
}

func MakeTarget(target string) task.Strategy {
	return &makeRun{Target: target}
}

// resolveTarget picks the first candidate from a comma-separated list
// that exists as a target in the Makefile. If only one candidate is given,
// it is used directly without probing.
func resolveTarget(dir, raw string) (string, error) {
	candidates := strings.Split(raw, ",")
	if len(candidates) == 1 {
		return strings.TrimSpace(candidates[0]), nil
	}

	for _, c := range candidates {
		c = strings.TrimSpace(c)
		if makeTargetExists(dir, c) {
			return c, nil
		}
	}
	return "", fmt.Errorf("none of the make targets [%s] found in %s", raw, dir)
}

func makeTargetExists(dir, target string) bool {
	cmd := exec.Command("make", "-n", target)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}
