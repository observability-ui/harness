package recipe

import (
	"fmt"
	"io/fs"
	"os/exec"

	"github.com/spf13/pflag"
)

type Status int

const (
	StatusPending Status = iota
	StatusWaiting // blocked on a dependency
	StatusRunning
	StatusStarted // processes launched and kept running (long-lived)
	StatusReady   // ports accepting connections (long-lived, confirmed ready)
	StatusDone    // processes completed and exited
	StatusStopped // processes intentionally stopped by user
	StatusFailed
	StatusSkipped
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusWaiting:
		return "waiting"
	case StatusRunning:
		return "running"
	case StatusStarted:
		return "started"
	case StatusReady:
		return "ready"
	case StatusDone:
		return "done"
	case StatusStopped:
		return "stopped"
	case StatusFailed:
		return "failed"
	case StatusSkipped:
		return "skipped"
	default:
		return "unknown"
	}
}

type Requirement struct {
	Name  string
	Check func() error
}

type FileRef struct {
	FS   fs.FS
	Path string
}

type ProcessSpec struct {
	Name      string
	Command   string
	Args      []string
	Dir       string
	Env       map[string]string
	Ports     []int
	Stdin     string
	StdinFile string
	Files     map[string]FileRef
}

type Step struct {
	Name      string
	Processes []ProcessSpec
	DependsOn []string
}

func (s *Step) HasPorts() bool {
	for _, p := range s.Processes {
		if len(p.Ports) > 0 {
			return true
		}
	}
	return false
}

type Config struct {
	Flags  *pflag.FlagSet
	DryRun bool
}

type StepUpdate struct {
	StepName string
	Status   Status
	Err      error
}

type Recipe interface {
	Name() string
	Aliases() []string
	Description() string
	Flags() *pflag.FlagSet
	Requirements(flags *pflag.FlagSet) []Requirement
	Steps(cfg *Config) ([]*Step, error)
}

type StepNeed struct {
	Provider string
	Config   map[string]string
}

type NeedfulRecipe interface {
	Recipe
	Needs() []StepNeed
}

type StepProvider interface {
	Name() string
	Provide(needs []StepNeed, cfg *Config) ([]*Step, error)
}

func RequireNode() Requirement   { return RequireTool("node", "install via nvm or brew") }
func RequireNPM() Requirement    { return RequireTool("npm", "install via nvm or brew") }
func RequireGo() Requirement     { return RequireTool("go", "") }
func RequirePodman() Requirement { return RequireTool("podman", "install via brew or dnf") }
func RequireJQ() Requirement     { return RequireTool("jq", "install via brew or dnf") }

func RequireTool(name, hint string) Requirement {
	return Requirement{
		Name: name,
		Check: func() error {
			if _, err := exec.LookPath(name); err != nil {
				if hint != "" {
					return fmt.Errorf("%s is not installed — %s", name, hint)
				}
				return fmt.Errorf("%s is not installed", name)
			}
			return nil
		},
	}
}

func RequireFlag(flags *pflag.FlagSet, name, usage string) Requirement {
	return Requirement{
		Name: fmt.Sprintf("--%s", name),
		Check: func() error {
			if flags == nil {
				return fmt.Errorf("--%s is required — %s", name, usage)
			}
			f := flags.Lookup(name)
			if f == nil || !f.Changed {
				return fmt.Errorf("--%s is required — %s", name, usage)
			}
			return nil
		},
	}
}

func RequireOCLogin() Requirement {
	return Requirement{
		Name: "oc (logged in)",
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
	}
}
