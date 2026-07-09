package component

import "io/fs"

type Status int

const (
	StatusPending Status = iota
	StatusWaiting        // blocked on a dependency
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

type FileRef struct {
	FS   fs.FS
	Path string
}

type Lifecycle int

const (
	LifecycleOneShot     Lifecycle = iota // runs to completion
	LifecycleLongRunning                  // stays running; readiness via port probing
)

func (l Lifecycle) String() string {
	switch l {
	case LifecycleOneShot:
		return "one-shot"
	case LifecycleLongRunning:
		return "long-running"
	default:
		return "unknown"
	}
}

type ProcessSpec struct {
	Name      string
	Command   string
	Args      []string
	Dir       string
	Env       map[string]string
	Ports []int
	Files map[string]FileRef
}

type Step struct {
	Name      string
	Lifecycle Lifecycle
	Processes []ProcessSpec
	DependsOn []string
}

type StepUpdate struct {
	StepName string
	Status   Status
	Err      error
}
