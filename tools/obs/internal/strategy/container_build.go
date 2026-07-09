package strategy

import (
	"context"
	"fmt"

	"obs/internal/component"
	"obs/internal/runcontext"
)

type ContainerBuild struct{}

func (s *ContainerBuild) Name() string        { return "container-build" }
func (s *ContainerBuild) Requires() []string { return []string{"podman", "bash"} }

func (s *ContainerBuild) Build(_ context.Context, comp *component.Component, rc *runcontext.RunContext) (*component.Step, error) {
	image := comp.Config["image"]
	if image == "" {
		image = rc.Get(comp.Name, "image")
	}
	dockerfile := comp.Config["dockerfile"]
	if dockerfile == "" {
		dockerfile = "Dockerfile"
	}
	platform := comp.Config["platform"]
	if platform == "" {
		platform = "linux/amd64"
	}

	script := fmt.Sprintf(
		"podman build -f %s --platform=%s -t %s . && podman push %s",
		dockerfile, platform, image, image,
	)

	return &component.Step{
		Name:      comp.Name,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{
			{
				Name:    comp.Name,
				Command: "bash",
				Args:    []string{"-c", script},
				Dir:     comp.Dir,
			},
		},
	}, nil
}
