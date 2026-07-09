package strategy

import (
	"context"
	"fmt"

	"obs/internal/component"
	"obs/internal/runcontext"
)

type ContainerBuild struct{}

func (s *ContainerBuild) Requires() []string { return []string{"podman", "bash"} }

func (s *ContainerBuild) Execute(_ context.Context, comp *component.Component, rc *runcontext.RunContext) (*component.Step, error) {
	image := comp.Config["image"]
	if image == "" {
		image = rc.Get(comp.Name, "image")
	}
	if image == "" {
		return nil, fmt.Errorf("component %q: no image specified (set Config[\"image\"] or provide --image flag)", comp.Name)
	}
	dockerfile := comp.Config["dockerfile"]
	if dockerfile == "" {
		dockerfile = "Dockerfile"
	}
	platform := comp.Config["platform"]
	if platform == "" {
		platform = "linux/amd64"
	}

	script := "podman build -f \"$DOCKERFILE\" --platform=\"$PLATFORM\" -t \"$IMAGE\" . && podman push \"$IMAGE\""

	return &component.Step{
		Name:      comp.Name,
		Lifecycle: component.LifecycleOneShot,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{
			{
				Name:    comp.Name,
				Command: "bash",
				Args:    []string{"-c", script},
				Dir:     comp.Dir,
				Env: map[string]string{
					"IMAGE":      image,
					"DOCKERFILE": dockerfile,
					"PLATFORM":   platform,
				},
			},
		},
	}, nil
}
