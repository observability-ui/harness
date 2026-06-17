package recipe

import "github.com/spf13/pflag"

type Status int

const (
	StatusPending Status = iota
	StatusRunning
	StatusDone
	StatusFailed
	StatusSkipped
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusRunning:
		return "running"
	case StatusDone:
		return "done"
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

type ProcessSpec struct {
	Name    string
	Command string
	Args    []string
	Dir     string
	Env     map[string]string
	Ports   []int
}

type Step struct {
	Name      string
	Processes []ProcessSpec
	DependsOn []string
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
	Requirements() []Requirement
	Steps(cfg *Config) ([]*Step, error)
}
