package ui

import (
	"testing"

	"obs/internal/process"
	"obs/internal/recipe"
)

func newTestMainTab() MainTab {
	mt := NewMainTab()
	mt.AddStepWithProcesses("step1", []string{"proc-a", "proc-b"})
	mt.AddStepWithProcesses("step2", []string{"proc-c"})
	return mt
}

func TestUpdateStep_PropagatesDone(t *testing.T) {
	mt := newTestMainTab()
	mt.UpdateStep("step1", recipe.StatusDone, nil)

	step := mt.GetStep("step1")
	if step.Status != recipe.StatusDone {
		t.Fatalf("step status: got %v, want Done", step.Status)
	}
	for _, p := range step.Processes {
		if p.Status != recipe.StatusDone {
			t.Errorf("process %q: got %v, want Done", p.Name, p.Status)
		}
	}
}

func TestUpdateStep_PropagatesStopped(t *testing.T) {
	mt := newTestMainTab()
	mt.UpdateStep("step1", recipe.StatusStopped, nil)

	for _, p := range mt.GetStep("step1").Processes {
		if p.Status != recipe.StatusStopped {
			t.Errorf("process %q: got %v, want Stopped", p.Name, p.Status)
		}
	}
}

func TestUpdateStep_PropagatesFailed(t *testing.T) {
	mt := newTestMainTab()
	mt.UpdateStep("step1", recipe.StatusFailed, nil)

	for _, p := range mt.GetStep("step1").Processes {
		if p.Status != recipe.StatusFailed {
			t.Errorf("process %q: got %v, want Failed", p.Name, p.Status)
		}
	}
}

func TestUpdateStep_PropagatesReady(t *testing.T) {
	mt := newTestMainTab()
	mt.UpdateStep("step1", recipe.StatusReady, nil)

	for _, p := range mt.GetStep("step1").Processes {
		if p.Status != recipe.StatusReady {
			t.Errorf("process %q: got %v, want Ready", p.Name, p.Status)
		}
	}
}

func TestUpdateStep_DoesNotPropagateStarted(t *testing.T) {
	mt := newTestMainTab()
	mt.UpdateStep("step1", recipe.StatusStarted, nil)

	for _, p := range mt.GetStep("step1").Processes {
		if p.Status != recipe.StatusPending {
			t.Errorf("process %q: got %v, want Pending (unchanged)", p.Name, p.Status)
		}
	}
}

func TestUpdateStep_DoesNotAffectOtherSteps(t *testing.T) {
	mt := newTestMainTab()
	mt.UpdateStep("step1", recipe.StatusDone, nil)

	step2 := mt.GetStep("step2")
	if step2.Status != recipe.StatusPending {
		t.Errorf("step2 status: got %v, want Pending", step2.Status)
	}
	for _, p := range step2.Processes {
		if p.Status != recipe.StatusPending {
			t.Errorf("step2 process %q: got %v, want Pending", p.Name, p.Status)
		}
	}
}

func TestUpdateProcess_IndividualUpdate(t *testing.T) {
	mt := newTestMainTab()
	mt.UpdateProcess("step1", "proc-a", recipe.StatusStarted)

	step := mt.GetStep("step1")
	if step.Processes[0].Status != recipe.StatusStarted {
		t.Errorf("proc-a: got %v, want Started", step.Processes[0].Status)
	}
	if step.Processes[1].Status != recipe.StatusPending {
		t.Errorf("proc-b: got %v, want Pending (unchanged)", step.Processes[1].Status)
	}
}

func TestMapProcessStatus(t *testing.T) {
	tests := []struct {
		input process.ProcessStatus
		want  recipe.Status
	}{
		{process.ProcessFailed, recipe.StatusFailed},
		{process.ProcessStopped, recipe.StatusStopped},
		{process.ProcessDone, recipe.StatusDone},
		{process.ProcessPending, recipe.StatusDone},
	}
	for _, tt := range tests {
		got := MapProcessStatus(tt.input)
		if got != tt.want {
			t.Errorf("MapProcessStatus(%v): got %v, want %v", tt.input, got, tt.want)
		}
	}
}
