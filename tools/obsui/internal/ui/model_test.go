package ui_test

import (
	"testing"

	"obsui/internal/process"
	"obsui/internal/recipe"
	"obsui/internal/types"
	"obsui/internal/ui"
)

func TestModel_Init(t *testing.T) {
	mgr := process.NewManager()
	steps := []*recipe.Step{
		{Name: "step1", Processes: []recipe.ProcessSpec{{Name: "proc1"}}},
	}
	updates := make(chan types.StepUpdate)
	close(updates)

	model := ui.NewModel(mgr, steps, updates)
	cmd := model.Init()
	if cmd == nil {
		t.Fatal("Init should return a Cmd")
	}
}
