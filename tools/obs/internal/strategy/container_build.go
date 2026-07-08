package strategy

import (
	"context"

	"obs/internal/component"
	"obs/internal/runcontext"
)

type ContainerBuild struct{}

func (s *ContainerBuild) Name() string        { return "container-build" }
func (s *ContainerBuild) Requires() []string { return []string{"podman"} }

func (s *ContainerBuild) Build(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	image := comp.Config["image"]
	dockerfile := comp.Config["dockerfile"]
	if dockerfile == "" {
		dockerfile = "Dockerfile"
	}
	platform := comp.Config["platform"]
	if platform == "" {
		platform = "linux/amd64"
	}

	buildArgs := []string{"build", "-f", dockerfile, "--platform=" + platform, "-t", image}
	pushArgs := []string{"push", image}

	return &component.Step{
		Name:      comp.Name,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{
			{
				Name:    comp.Name + "-build",
				Command: "podman",
				Args:    buildArgs,
				Dir:     comp.Dir,
			},
			{
				Name:    comp.Name + "-push",
				Command: "podman",
				Args:    pushArgs,
				Dir:     comp.Dir,
			},
		},
	}, nil
}
