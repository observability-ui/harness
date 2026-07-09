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

type stubRunStrategy struct {
	name string
	step *component.Step
}

func (s *stubRunStrategy) Name() string        { return s.name }
func (s *stubRunStrategy) Requires() []string { return nil }
func (s *stubRunStrategy) Run(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
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

func setup() {
	strategy.RegisterSelector(func(comp *component.Component, mode string) (strategy.BuildStrategy, strategy.RunStrategy) {
		return nil, &stubRunStrategy{}
	})
}

func TestMix_SingleComponent(t *testing.T) {
	setup()
	component.Register(&component.Component{
		Name:        "test-a",
		Description: "Test component A",
	})

	m := mixer.New()
	steps, rc, err := m.Mix(context.Background(), []string{"test-a"}, "local", nil)
	if err != nil {
		t.Fatalf("Mix failed: %v", err)
	}
	if len(steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(steps))
	}
	if steps[0].Name != "test-a" {
		t.Fatalf("expected step name test-a, got %s", steps[0].Name)
	}
	if rc.Mode() != "local" {
		t.Fatalf("expected mode local, got %s", rc.Mode())
	}
}

func TestMix_DependencyOrder(t *testing.T) {
	setup()
	component.Register(&component.Component{
		Name: "test-dep",
	})
	component.Register(&component.Component{
		Name:      "test-main",
		DependsOn: []string{"test-dep"},
	})

	m := mixer.New()
	steps, _, err := m.Mix(context.Background(), []string{"test-main"}, "local", nil)
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
	setup()
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

	m := mixer.New()
	steps, _, err := m.Mix(context.Background(), []string{"app-a", "app-b"}, "local", nil)
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

func TestMix_OutputsWrittenToRunContext(t *testing.T) {
	setup()
	component.Register(&component.Component{
		Name:    "svc",
		Outputs: []component.Output{{Name: "port", Value: "8080"}},
	})

	m := mixer.New()
	_, rc, err := m.Mix(context.Background(), []string{"svc"}, "local", nil)
	if err != nil {
		t.Fatalf("Mix failed: %v", err)
	}
	if rc.Get("svc", "port") != "8080" {
		t.Fatalf("expected port 8080 in RunContext, got %q", rc.Get("svc", "port"))
	}
}

func TestMix_CircularDependency(t *testing.T) {
	setup()
	component.Register(&component.Component{
		Name:      "cycle-a",
		DependsOn: []string{"cycle-b"},
	})
	component.Register(&component.Component{
		Name:      "cycle-b",
		DependsOn: []string{"cycle-a"},
	})

	m := mixer.New()
	_, _, err := m.Mix(context.Background(), []string{"cycle-a"}, "local", nil)
	if err == nil {
		t.Fatal("expected error for circular dependency")
	}
	if !strings.Contains(err.Error(), "circular dependency") {
		t.Fatalf("expected circular dependency error, got: %v", err)
	}
}

func TestMix_UnknownComponent(t *testing.T) {
	m := mixer.New()
	_, _, err := m.Mix(context.Background(), []string{"nonexistent"}, "local", nil)
	if err == nil {
		t.Fatal("expected error for unknown component")
	}
}

func TestMix_RequiredFlagsPropagatedToComponentScope(t *testing.T) {
	setup()
	component.Register(&component.Component{
		Name: "flagged-comp",
		RequiredFlags: []component.RequiredFlag{
			{Name: "image", Usage: "container image"},
		},
	})

	m := mixer.New()
	flagValues := map[string]string{"image": "quay.io/test/img:v1"}
	_, rc, err := m.Mix(context.Background(), []string{"flagged-comp"}, "local", flagValues)
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
	setup()
	component.Register(&component.Component{
		Name: "needs-flag",
		RequiredFlags: []component.RequiredFlag{
			{Name: "token", Usage: "auth token"},
		},
	})

	m := mixer.New()
	_, _, err := m.Mix(context.Background(), []string{"needs-flag"}, "local", nil)
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
