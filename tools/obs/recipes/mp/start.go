package mp

import (
	"obs/internal/recipe"

	"github.com/spf13/pflag"
)

type StartMonitoringPlugin struct{}

func (r *StartMonitoringPlugin) Name() string      { return "monitoring-plugin" }
func (r *StartMonitoringPlugin) Aliases() []string { return []string{"mp"} }
func (r *StartMonitoringPlugin) Description() string {
	return "Start monitoring plugin frontend and backend dev servers"
}

func (r *StartMonitoringPlugin) Flags() *pflag.FlagSet {
	return pflag.NewFlagSet("monitoring-plugin", pflag.ContinueOnError)
}

func (r *StartMonitoringPlugin) Requirements() []recipe.Requirement {
	return []recipe.Requirement{
		recipe.RequireNode(),
		recipe.RequireNPM(),
		recipe.RequireGo(),
		recipe.RequireOCLogin(),
	}
}

func (r *StartMonitoringPlugin) Steps(cfg *recipe.Config) ([]*recipe.Step, error) {
	dir := "projects/monitoring-plugin"

	return []*recipe.Step{
		{
			Name:      "start-mp-frontend",
			DependsOn: []string{},
			Processes: []recipe.ProcessSpec{
				{
					Name:    "mp-frontend",
					Command: "make",
					Args:    []string{"start-frontend"},
					Dir:     dir,
					Ports:   []int{9001},
				},
			},
		},
		{
			Name:      "start-mp-backend",
			DependsOn: []string{},
			Processes: []recipe.ProcessSpec{
				{
					Name:    "mp-backend",
					Command: "make",
					Args:    []string{"start-feature-backend"},
					Dir:     dir,
					Ports:   []int{9443},
				},
			},
		},
		{
			Name:      "start-mp-console",
			DependsOn: []string{},
			Processes: []recipe.ProcessSpec{
				{
					Name:    "mp-console",
					Command: "make",
					Args:    []string{"start-feature-console"},
					Dir:     dir,
					Ports:   []int{9000},
				},
			},
		},
	}, nil
}
