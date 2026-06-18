package runner

import (
	"context"

	"obs/internal/process"
	"obs/internal/recipe"
)

type Runner interface {
	Run(ctx context.Context, mgr *process.Manager, steps []*recipe.Step, updates chan<- recipe.StepUpdate) error
}
