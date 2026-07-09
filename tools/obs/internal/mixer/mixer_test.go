package mixer_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"obs/internal/component"
	"obs/internal/mixer"
	"obs/internal/runcontext"
	"obs/internal/strategy"
)

type stubStrategy struct {
	step *component.Step
}

func (s *stubStrategy) Requires() []string { return nil }
func (s *stubStrategy) Execute(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	if s.step != nil {
		return s.step, nil
	}
	return &component.Step{
		Name:      comp.Name,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{{
			Name:    comp.Name,
			Command: "echo",
			Args:    []string{comp.Name},
		}},
	}, nil
}

func setup(t *testing.T, componentNames ...string) {
	t.Helper()
	t.Cleanup(func() {
		component.ResetRegistry()
		strategy.ResetRegistry()
		mixer.ResetRecipes()
	})
	for _, name := range componentNames {
		strategy.Register(name, &stubStrategy{})
	}
}

func TestMix_SingleComponent(t *testing.T) {
	setup(t, "test-a")
	component.Register(&component.Component{
		Name: "test-a",
	})

	steps, _, err := mixer.Mix(context.Background(), []string{"test-a"}, nil)
	if err != nil {
		t.Fatalf("Mix failed: %v", err)
	}
	if len(steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(steps))
	}
	if steps[0].Name != "test-a" {
		t.Fatalf("expected step name test-a, got %s", steps[0].Name)
	}
}

func TestMix_DependencyOrder(t *testing.T) {
	setup(t, "test-dep", "test-main")
	component.Register(&component.Component{
		Name: "test-dep",
	})
	component.Register(&component.Component{
		Name:      "test-main",
		DependsOn: []string{"test-dep"},
	})

	steps, _, err := mixer.Mix(context.Background(), []string{"test-main"}, nil)
	if err != nil {
		t.Fatalf("Mix failed: %v", err)
	}
	if len(steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(steps))
	}
	if steps[0].Name != "test-dep" {
		t.Fatalf("expected test-dep first, got %s", steps[0].Name)
	}
}

func TestMix_DeduplicatesSharedComponents(t *testing.T) {
	setup(t, "shared-infra", "app-a", "app-b")
	component.Register(&component.Component{
		Name: "shared-infra",
	})
	component.Register(&component.Component{
		Name:      "app-a",
		DependsOn: []string{"shared-infra"},
	})
	component.Register(&component.Component{
		Name:      "app-b",
		DependsOn: []string{"shared-infra"},
	})

	steps, _, err := mixer.Mix(context.Background(), []string{"app-a", "app-b"}, nil)
	if err != nil {
		t.Fatalf("Mix failed: %v", err)
	}

	count := 0
	for _, s := range steps {
		if s.Name == "shared-infra" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected shared-infra once, got %d", count)
	}
}

func TestMix_PortsWrittenToRunContext(t *testing.T) {
	setup(t, "svc")
	component.Register(&component.Component{
		Name:  "svc",
		Ports: []int{8080},
	})

	_, rc, err := mixer.Mix(context.Background(), []string{"svc"}, nil)
	if err != nil {
		t.Fatalf("Mix failed: %v", err)
	}
	if rc.Get("svc", "port") != "8080" {
		t.Fatalf("expected port 8080 in RunContext, got %q", rc.Get("svc", "port"))
	}
}

func TestMix_CircularDependency(t *testing.T) {
	setup(t, "cycle-a", "cycle-b")
	component.Register(&component.Component{
		Name:      "cycle-a",
		DependsOn: []string{"cycle-b"},
	})
	component.Register(&component.Component{
		Name:      "cycle-b",
		DependsOn: []string{"cycle-a"},
	})

	_, _, err := mixer.Mix(context.Background(), []string{"cycle-a"}, nil)
	if err == nil {
		t.Fatal("expected error for circular dependency")
	}
	if !strings.Contains(err.Error(), "circular dependency") {
		t.Fatalf("expected circular dependency error, got: %v", err)
	}
}

func TestMix_UnknownComponent(t *testing.T) {
	_, _, err := mixer.Mix(context.Background(), []string{"nonexistent"}, nil)
	if err == nil {
		t.Fatal("expected error for unknown component")
	}
}

func TestMix_RequiredFlagsPropagatedToComponentScope(t *testing.T) {
	setup(t, "flagged-comp")
	component.Register(&component.Component{
		Name: "flagged-comp",
		RequiredFlags: []component.RequiredFlag{
			{Name: "image", Usage: "container image"},
		},
	})

	flagValues := map[string]string{"image": "quay.io/test/img:v1"}
	_, rc, err := mixer.Mix(context.Background(), []string{"flagged-comp"}, flagValues)
	if err != nil {
		t.Fatalf("Mix failed: %v", err)
	}

	if v := rc.Get("_flags", "image"); v != "quay.io/test/img:v1" {
		t.Errorf("expected flag under _flags, got %q", v)
	}
	if v := rc.Get("flagged-comp", "image"); v != "quay.io/test/img:v1" {
		t.Errorf("expected flag under component name, got %q", v)
	}
}

func TestMix_MissingRequiredFlagReturnsError(t *testing.T) {
	setup(t, "needs-flag")
	component.Register(&component.Component{
		Name: "needs-flag",
		RequiredFlags: []component.RequiredFlag{
			{Name: "token", Usage: "auth token"},
		},
	})

	_, _, err := mixer.Mix(context.Background(), []string{"needs-flag"}, nil)
	if err == nil {
		t.Fatal("expected error for missing required flag")
	}
	var mfe *mixer.MissingFlagError
	if !errors.As(err, &mfe) {
		t.Fatalf("expected MissingFlagError, got %T: %v", err, err)
	}
	if mfe.Flag != "token" {
		t.Errorf("expected flag name 'token', got %q", mfe.Flag)
	}
}
