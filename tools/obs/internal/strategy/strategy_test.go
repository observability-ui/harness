package strategy

import (
	"context"
	"testing"

	"obs/internal/component"
	"obs/internal/runcontext"
)

type testStrategy struct{ label string }

func (s *testStrategy) Requires() []string { return nil }
func (s *testStrategy) Execute(_ context.Context, _ *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	return nil, nil
}

func TestResolve_ReturnsNilForUnknownComponent(t *testing.T) {
	t.Cleanup(func() { ResetRegistry() })

	comp := &component.Component{Name: "unknown"}
	strategies := Resolve(comp)
	if strategies != nil {
		t.Fatalf("expected nil for unknown component, got %v", strategies)
	}
}

func TestResolve_ExplicitRegistration(t *testing.T) {
	t.Cleanup(func() { ResetRegistry() })

	s := &testStrategy{label: "explicit"}
	Register("my-comp", s)

	comp := &component.Component{Name: "my-comp"}
	strategies := Resolve(comp)
	if len(strategies) != 1 || strategies[0] != s {
		t.Fatalf("expected explicitly registered strategy, got %v", strategies)
	}
}

func TestResolve_FallsBackToConfig(t *testing.T) {
	t.Cleanup(func() { ResetRegistry() })

	comp := &component.Component{
		Name:   "unregistered",
		Config: map[string]string{"make-target": "build"},
	}
	strategies := Resolve(comp)
	if len(strategies) != 1 {
		t.Fatalf("expected config-based fallback to return 1 strategy, got %d", len(strategies))
	}
}
