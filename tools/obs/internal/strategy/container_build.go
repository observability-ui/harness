package strategy

import (
	"context"
	"fmt"

	"obs/internal/runcontext"
	"obs/internal/task"
)

type containerBuild struct {
	Dockerfile string
	Platform   string
}

func (s *containerBuild) Requires() []string { return []string{"podman", "bash"} }

func (s *containerBuild) Execute(_ context.Context, t *task.Task, rc *runcontext.RunContext) (*task.Step, error) {
	image := rc.Get(t.Name, "image")
	if image == "" {
		return nil, fmt.Errorf("task %q: no image specified (provide --image flag)", t.Name)
	}

	dockerfile := s.Dockerfile
	if dockerfile == "" {
		dockerfile = "Dockerfile"
	}
	platform := s.Platform
	if platform == "" {
		platform = "linux/amd64"
	}

	script := "podman build -f \"$DOCKERFILE\" --platform=\"$PLATFORM\" -t \"$IMAGE\" . && podman push \"$IMAGE\""

	return &task.Step{
		Name:      t.Name,
		Lifecycle: task.LifecycleOneShot,
		DependsOn: t.DependsOn,
		Processes: []task.ProcessSpec{
			{
				Name:    t.Name,
				Command: "bash",
				Args:    []string{"-c", script},
				Dir:     t.Dir,
				Env: map[string]string{
					"IMAGE":      image,
					"DOCKERFILE": dockerfile,
					"PLATFORM":   platform,
				},
			},
		},
	}, nil
}

func DockerBuild(dockerfile string) task.Strategy {
	return &containerBuild{Dockerfile: dockerfile}
}
