package mixer_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"obs/internal/mixer"
	"obs/internal/runcontext"
	"obs/internal/task"
)

type stubStrategy struct {
	step *task.Step
}

func (s *stubStrategy) Requires() []string { return nil }
func (s *stubStrategy) Execute(_ context.Context, t *task.Task, _ *runcontext.RunContext) (*task.Step, error) {
	if s.step != nil {
		return s.step, nil
	}
	return &task.Step{
		Name:      t.Name,
		DependsOn: t.DependsOn,
		Processes: []task.ProcessSpec{{
			Name:    t.Name,
			Command: "echo",
			Args:    []string{t.Name},
		}},
	}, nil
}

func setup(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		task.ResetRegistry()
		mixer.ResetRecipes()
	})
}

func TestMix_SingleTask(t *testing.T) {
	setup(t)
	task.Register(&task.Task{
		Name:     "test-a",
		Strategy: &stubStrategy{},
	})

	steps, _, _, err := mixer.Mix(context.Background(), []string{"test-a"}, nil)
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
	setup(t)
	task.Register(&task.Task{
		Name:     "test-dep",
		Strategy: &stubStrategy{},
	})
	task.Register(&task.Task{
		Name:      "test-main",
		DependsOn: []string{"test-dep"},
		Strategy:  &stubStrategy{},
	})

	steps, _, _, err := mixer.Mix(context.Background(), []string{"test-main"}, nil)
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

func TestMix_DeduplicatesSharedTasks(t *testing.T) {
	setup(t)
	task.Register(&task.Task{
		Name:     "shared-infra",
		Strategy: &stubStrategy{},
	})
	task.Register(&task.Task{
		Name:      "app-a",
		DependsOn: []string{"shared-infra"},
		Strategy:  &stubStrategy{},
	})
	task.Register(&task.Task{
		Name:      "app-b",
		DependsOn: []string{"shared-infra"},
		Strategy:  &stubStrategy{},
	})

	steps, _, _, err := mixer.Mix(context.Background(), []string{"app-a", "app-b"}, nil)
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
	setup(t)
	task.Register(&task.Task{
		Name:     "svc",
		Ports:    []int{8080},
		Strategy: &stubStrategy{},
	})

	_, rc, _, err := mixer.Mix(context.Background(), []string{"svc"}, nil)
	if err != nil {
		t.Fatalf("Mix failed: %v", err)
	}
	if rc.Get("svc", "port") != "8080" {
		t.Fatalf("expected port 8080 in RunContext, got %q", rc.Get("svc", "port"))
	}
}

func TestMix_CircularDependency(t *testing.T) {
	setup(t)
	task.Register(&task.Task{
		Name:      "cycle-a",
		DependsOn: []string{"cycle-b"},
		Strategy:  &stubStrategy{},
	})
	task.Register(&task.Task{
		Name:      "cycle-b",
		DependsOn: []string{"cycle-a"},
		Strategy:  &stubStrategy{},
	})

	_, _, _, err := mixer.Mix(context.Background(), []string{"cycle-a"}, nil)
	if err == nil {
		t.Fatal("expected error for circular dependency")
	}
	if !strings.Contains(err.Error(), "circular dependency") {
		t.Fatalf("expected circular dependency error, got: %v", err)
	}
}

func TestMix_UnknownTask(t *testing.T) {
	setup(t)
	_, _, _, err := mixer.Mix(context.Background(), []string{"nonexistent"}, nil)
	if err == nil {
		t.Fatal("expected error for unknown task")
	}
}

func TestMix_RequiredFlagsPropagatedToTaskScope(t *testing.T) {
	setup(t)
	task.Register(&task.Task{
		Name:     "flagged-task",
		Strategy: &stubStrategy{},
		RequiredFlags: []task.RequiredFlag{
			{Name: "image", Usage: "container image"},
		},
	})

	flagValues := map[string]string{"image": "quay.io/test/img:v1"}
	_, rc, _, err := mixer.Mix(context.Background(), []string{"flagged-task"}, flagValues)
	if err != nil {
		t.Fatalf("Mix failed: %v", err)
	}

	if v := rc.Get("_flags", "image"); v != "quay.io/test/img:v1" {
		t.Errorf("expected flag under _flags, got %q", v)
	}
	if v := rc.Get("flagged-task", "image"); v != "quay.io/test/img:v1" {
		t.Errorf("expected flag under task name, got %q", v)
	}
}

func TestMix_MissingRequiredFlagReturnsError(t *testing.T) {
	setup(t)
	task.Register(&task.Task{
		Name:     "needs-flag",
		Strategy: &stubStrategy{},
		RequiredFlags: []task.RequiredFlag{
			{Name: "token", Usage: "auth token"},
		},
	})

	_, _, _, err := mixer.Mix(context.Background(), []string{"needs-flag"}, nil)
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
