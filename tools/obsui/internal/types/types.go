package types

import "obsui/internal/recipe"

// StepUpdate is sent by runners to report step progress
type StepUpdate struct {
	StepName string
	Status   recipe.Status
	Err      error
}
