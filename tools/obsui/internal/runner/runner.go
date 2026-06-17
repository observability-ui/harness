package runner

import (
	"context"

	"obsui/internal/process"
	"obsui/internal/recipe"
	"obsui/internal/types"
)

// StepUpdate is an alias for types.StepUpdate for backward compatibility
type StepUpdate = types.StepUpdate

type Runner interface {
	Run(ctx context.Context, mgr *process.Manager, steps []*recipe.Step, updates chan<- types.StepUpdate) error
}
