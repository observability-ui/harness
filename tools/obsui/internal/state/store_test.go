package state_test

import (
	"os"
	"path/filepath"
	"testing"

	"obsui/internal/state"
)

func TestStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store := state.NewStore(dir)

	s := &state.RunState{
		Recipes: []string{"monitoring-plugin"},
		Processes: []state.ProcessState{
			{Name: "mp-frontend", PID: 12345, Status: "running", Port: 9000},
			{Name: "mp-backend", PID: 12346, Status: "running", Port: 9001},
		},
	}

	if err := store.Save(s); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Processes) != 2 {
		t.Fatalf("expected 2 processes, got %d", len(loaded.Processes))
	}
	if loaded.Processes[0].PID != 12345 {
		t.Fatalf("PID mismatch: %d", loaded.Processes[0].PID)
	}
}

func TestStore_Clean(t *testing.T) {
	dir := t.TempDir()
	store := state.NewStore(dir)

	store.Save(&state.RunState{Processes: []state.ProcessState{{Name: "test", PID: 1}}})
	store.Clean()

	if _, err := os.Stat(filepath.Join(dir, "state.json")); !os.IsNotExist(err) {
		t.Fatal("state.json should be removed after Clean")
	}
}

func TestStore_WritePID(t *testing.T) {
	dir := t.TempDir()
	store := state.NewStore(dir)

	if err := store.WritePID("test-proc", 99999); err != nil {
		t.Fatalf("WritePID failed: %v", err)
	}

	pid, err := store.ReadPID("test-proc")
	if err != nil {
		t.Fatalf("ReadPID failed: %v", err)
	}
	if pid != 99999 {
		t.Fatalf("expected PID 99999, got %d", pid)
	}
}
