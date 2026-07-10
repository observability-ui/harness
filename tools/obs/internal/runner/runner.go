package runner

import (
	"context"

	"obs/internal/task"
	"obs/internal/process"
)

type Runner interface {
	Run(ctx context.Context, mgr *process.Manager, steps []*task.Step, updates chan<- task.StepUpdate) error
}
