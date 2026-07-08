package runner

import (
	"context"

	"obs/internal/component"
	"obs/internal/process"
)

type Runner interface {
	Run(ctx context.Context, mgr *process.Manager, steps []*component.Step, updates chan<- component.StepUpdate) error
}
