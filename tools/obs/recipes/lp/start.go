package lp

import (
	"obs/internal/recipe"

	"github.com/spf13/pflag"
)

type StartLoggingPlugin struct{}

func (r *StartLoggingPlugin) Name() string      { return "logging-view-plugin" }
func (r *StartLoggingPlugin) Aliases() []string { return []string{"lp"} }
func (r *StartLoggingPlugin) Description() string {
	return "Start logging view plugin frontend and backend dev servers"
}

func (r *StartLoggingPlugin) Flags() *pflag.FlagSet {
	return pflag.NewFlagSet("logging-view-plugin", pflag.ContinueOnError)
}

func (r *StartLoggingPlugin) Requirements(_ *pflag.FlagSet) []recipe.Requirement {
	return []recipe.Requirement{
		recipe.RequireNode(),
		recipe.RequireNPM(),
		recipe.RequireGo(),
		recipe.RequireOCLogin(),
		recipe.RequirePodman(),
	}
}

func (r *StartLoggingPlugin) Steps(cfg *recipe.Config) ([]*recipe.Step, error) {
	dir := "projects/logging-view-plugin"

	return []*recipe.Step{
		{
			Name:      "start-lp-frontend",
			DependsOn: []string{},
			Processes: []recipe.ProcessSpec{
				{
					Name:    "start lp frontend",
					Command: "make",
					Args:    []string{"start-frontend"},
					Dir:     dir,
					Ports:   []int{9001},
				},
			},
		},
		{
			Name:      "start-lp-backend",
			DependsOn: []string{"start-lp-frontend"},
			Processes: []recipe.ProcessSpec{
				{
					Name:    "start lp backend",
					Command: "make",
					Args:    []string{"start-backend"},
					Dir:     dir,
					Ports:   []int{9002},
				},
			},
		},
		{
			Name:      "start-lp-console",
			DependsOn: []string{},
			Processes: []recipe.ProcessSpec{
				{
					Name:    "start lp console",
					Command: "make",
					Args:    []string{"start-console"},
					Dir:     dir,
					Ports:   []int{9000},
				},
			},
		},
		{
			Name:      "start-lp-local-loki",
			DependsOn: []string{},
			Processes: []recipe.ProcessSpec{
				{
					Name:    "lp-local-loki",
					Command: "podman",
					Args:    []string{"compose", "-f", "hack/docker-compose/docker-compose.test.yml", "up"},
					Dir:     dir,
					Ports:   []int{3100},
				},
			},
		},
	}, nil
}
