package startmp

import (
	"fmt"
	"os/exec"

	"github.com/spf13/pflag"
	"obsui/internal/recipe"
)

type StartMonitoringPlugin struct{}

func (r *StartMonitoringPlugin) Name() string        { return "monitoring-plugin" }
func (r *StartMonitoringPlugin) Aliases() []string    { return []string{"mp"} }
func (r *StartMonitoringPlugin) Description() string  { return "Start monitoring plugin frontend and backend dev servers" }

func (r *StartMonitoringPlugin) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("monitoring-plugin", pflag.ContinueOnError)
	fs.String("version", "", "OpenShift version branch to use (e.g., 4.18)")
	return fs
}

func (r *StartMonitoringPlugin) Requirements() []recipe.Requirement {
	return []recipe.Requirement{
		{
			Name: "node",
			Check: func() error {
				if _, err := exec.LookPath("node"); err != nil {
					return fmt.Errorf("node is not installed — install via nvm or brew")
				}
				return nil
			},
		},
		{
			Name: "go",
			Check: func() error {
				if _, err := exec.LookPath("go"); err != nil {
					return fmt.Errorf("go is not installed")
				}
				return nil
			},
		},
	}
}

func (r *StartMonitoringPlugin) Steps(cfg *recipe.Config) ([]*recipe.Step, error) {
	version, _ := cfg.Flags.GetString("version")
	dir := "projects/monitoring-plugin"
	if version != "" {
		_ = version // branch checkout would happen separately
	}

	return []*recipe.Step{
		{
			Name: "start-mp-backend",
			Processes: []recipe.ProcessSpec{
				{
					Name:    "mp-backend",
					Command: "make",
					Args:    []string{"run-backend"},
					Dir:     dir,
					Ports:   []int{9443},
				},
			},
		},
		{
			Name: "start-mp-frontend",
			Processes: []recipe.ProcessSpec{
				{
					Name:    "mp-frontend",
					Command: "make",
					Args:    []string{"run-frontend"},
					Dir:     dir,
					Ports:   []int{9001},
				},
			},
		},
	}, nil
}
