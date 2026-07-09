package process_test

import (
	"context"
	"testing"
	"time"

	"obs/internal/component"
	"obs/internal/process"
)

func TestManager_StartAndStop(t *testing.T) {
	mgr := process.NewManager()
	ctx := context.Background()

	spec := component.ProcessSpec{
		Name:    "sleeper",
		Command: "sleep",
		Args:    []string{"30"},
	}

	proc, err := mgr.StartProcess(ctx, spec)
	if err != nil {
		t.Fatalf("StartProcess failed: %v", err)
	}
	if !proc.Running() {
		t.Fatal("process should be running")
	}

	// Duplicate start should fail
	_, err = mgr.StartProcess(ctx, spec)
	if err == nil {
		t.Fatal("duplicate start should fail")
	}

	mgr.StopAll()

	select {
	case <-proc.Wait():
	case <-time.After(10 * time.Second):
		t.Fatal("process did not stop in time")
	}
}

func TestManager_OutputCapture(t *testing.T) {
	mgr := process.NewManager()
	ctx := context.Background()

	spec := component.ProcessSpec{
		Name:    "echo",
		Command: "echo",
		Args:    []string{"hello world"},
	}

	proc, err := mgr.StartProcess(ctx, spec)
	if err != nil {
		t.Fatalf("StartProcess failed: %v", err)
	}
	<-proc.Wait()

	// Give the ring buffer a moment to flush
	time.Sleep(50 * time.Millisecond)

	lines := proc.Output.Lines()
	if len(lines) == 0 {
		t.Fatal("expected captured output")
	}
	if lines[0] != "hello world" {
		t.Fatalf("expected 'hello world', got %q", lines[0])
	}
}
