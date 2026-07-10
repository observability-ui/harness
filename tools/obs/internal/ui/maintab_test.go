package ui

import (
	"testing"

	"obs/internal/task"
)

func newTestMainTab() mainTab {
	mt := newMainTab()
	mt.AddStepWithProcesses("step1", []string{"proc-a", "proc-b"}, nil)
	mt.AddStepWithProcesses("step2", []string{"proc-c"}, []string{"step1"})
	return mt
}

func TestUpdateStep_SetsStepStatus(t *testing.T) {
	mt := newTestMainTab()
	mt.UpdateStep("step1", task.StatusDone, nil)

	step := mt.GetStep("step1")
	if step.Status != task.StatusDone {
		t.Fatalf("step status: got %v, want Done", step.Status)
	}
	for _, p := range step.Processes {
		if p.Status != task.StatusPending {
			t.Errorf("process %q: got %v, want Pending (UpdateStep does not change per-process status)", p.Name, p.Status)
		}
	}
}

func TestUpdateStep_DoesNotAffectOtherSteps(t *testing.T) {
	mt := newTestMainTab()
	mt.UpdateStep("step1", task.StatusDone, nil)

	step2 := mt.GetStep("step2")
	if step2.Status != task.StatusPending {
		t.Errorf("step2 status: got %v, want Pending", step2.Status)
	}
}

func TestUpdateProcess_IndividualUpdate(t *testing.T) {
	mt := newTestMainTab()
	mt.updateProcess("step1", "proc-a", task.StatusStarted)

	step := mt.GetStep("step1")
	if step.Processes[0].Status != task.StatusStarted {
		t.Errorf("proc-a: got %v, want Started", step.Processes[0].Status)
	}
	if step.Processes[1].Status != task.StatusPending {
		t.Errorf("proc-b: got %v, want Pending (unchanged)", step.Processes[1].Status)
	}
}

