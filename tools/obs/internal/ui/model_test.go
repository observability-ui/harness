package ui_test

import (
	"context"
	"testing"

	"obs/internal/task"
	"obs/internal/process"
	"obs/internal/ui"
)

func TestModel_Init(t *testing.T) {
	mgr := process.NewManager()
	updates := make(chan task.StepUpdate)
	close(updates)

	retryCh := make(chan struct{}, 1)
	model := ui.NewModel(context.Background(), mgr, updates, retryCh, nil)
	cmd := model.Init()
	if cmd == nil {
		t.Fatal("Init should return a Cmd")
	}
}
