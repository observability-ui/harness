package runner

import (
	"context"

	"obsui/internal/process"
	"obsui/internal/recipe"
)

type StepUpdate struct {
	StepName string
	Status   recipe.Status
	Err      error
}

type Runner interface {
	Run(ctx context.Context, mgr *process.Manager, steps []*recipe.Step, updates chan<- StepUpdate) error
}
