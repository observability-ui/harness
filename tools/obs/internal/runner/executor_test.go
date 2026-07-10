package runner

import (
	"context"
	"sync"
	"testing"
	"time"

	"obs/internal/task"
	"obs/internal/process"
)

func TestExecuteSteps_RunsStepsInOrder(t *testing.T) {
	ctx := context.Background()
	mgr := process.NewManager()
	defer mgr.StopAll()

	var order []string
	var mu sync.Mutex

	steps := []*task.Step{
		{Name: "step-a", Processes: []task.ProcessSpec{{Name: "a", Command: "echo", Args: []string{"a"}}}},
		{Name: "step-b", DependsOn: []string{"step-a"}, Processes: []task.ProcessSpec{{Name: "b", Command: "echo", Args: []string{"b"}}}},
	}

	cb := StepCallbacks{
		OnUpdate: func(u task.StepUpdate) {
			if u.Status == task.StatusStarted {
				mu.Lock()
				order = append(order, u.StepName)
				mu.Unlock()
			}
		},
		OnProcess: func(_ *task.Step, _ task.ProcessSpec, _ *process.Process) {},
	}

	_, err := ExecuteSteps(ctx, mgr, steps, cb)
	if err != nil {
		t.Fatalf("ExecuteSteps failed: %v", err)
	}

	// Give processes time to complete
	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(order) < 2 {
		t.Fatalf("expected 2 started steps, got %d", len(order))
	}
	if order[0] != "step-a" {
		t.Errorf("expected step-a first, got %s", order[0])
	}
}
