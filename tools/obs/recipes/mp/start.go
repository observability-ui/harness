package mp

import (
	"fmt"
	"os/exec"

	"obs/internal/process"
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
			Name: "npm",
			Check: func() error {
				if _, err := exec.LookPath("npm"); err != nil {
					return fmt.Errorf("npm is not installed — install via nvm or brew")
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
		{
			Name: "oc",
			Check: func() error {
				if _, err := exec.LookPath("oc"); err != nil {
					return fmt.Errorf("oc is not installed — install the OpenShift CLI")
				}
				out, err := exec.Command("oc", "whoami").CombinedOutput()
				if err != nil {
					return fmt.Errorf("not logged in to OpenShift cluster — run 'oc login' first (oc whoami: %s)", out)
				}
				return nil
			},
		},
		{
			Name:  "port 9001 (frontend)",
			Check: func() error { return process.CheckPort(9001) },
		},
		{
			Name:  "port 9443 (backend)",
			Check: func() error { return process.CheckPort(9443) },
		},
		{
			Name:  "port 9000 (console)",
			Check: func() error { return process.CheckPort(9000) },
		},
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
