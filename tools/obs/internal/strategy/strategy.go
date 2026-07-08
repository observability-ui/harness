package strategy

import (
	"context"

	"obs/internal/component"
	"obs/internal/runcontext"
)

type BuildStrategy interface {
	Name() string
	Requires() []string
	Build(ctx context.Context, comp *component.Component, rc *runcontext.RunContext) (*component.Step, error)
}

type RunStrategy interface {
	Name() string
	Requires() []string
	Run(ctx context.Context, comp *component.Component, rc *runcontext.RunContext) (*component.Step, error)
}

type Selector func(comp *component.Component, mode string) (BuildStrategy, RunStrategy)

var selectors []Selector

func RegisterSelector(sel Selector) {
	selectors = append(selectors, sel)
}

func Select(comp *component.Component, mode string) (BuildStrategy, RunStrategy) {
	for _, sel := range selectors {
		build, run := sel(comp, mode)
		if build != nil || run != nil {
			return build, run
		}
	}
	return nil, nil
}
