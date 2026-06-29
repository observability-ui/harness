package ui_test

import (
	"testing"

	"obs/internal/process"
	"obs/internal/recipe"
	"obs/internal/ui"
)

func TestModel_Init(t *testing.T) {
	mgr := process.NewManager()
	updates := make(chan recipe.StepUpdate)
	close(updates)

	retryCh := make(chan struct{}, 1)
	model := ui.NewModel(mgr, updates, retryCh)
	cmd := model.Init()
	if cmd == nil {
		t.Fatal("Init should return a Cmd")
	}
}
