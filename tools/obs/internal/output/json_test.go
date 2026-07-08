package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"obs/internal/output"
	"obs/internal/component"
)

func TestJSONEmitter(t *testing.T) {
	var buf bytes.Buffer
	emitter := output.NewJSONEmitter(&buf)

	emitter.Emit(output.Event{
		Type:   "step_status",
		Step:   "install-clo",
		Status: component.StatusRunning.String(),
	})
	emitter.Emit(output.Event{
		Type:   "step_status",
		Step:   "install-clo",
		Status: component.StatusDone.String(),
	})

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 JSON lines, got %d", len(lines))
	}

	var ev output.Event
	if err := json.Unmarshal([]byte(lines[0]), &ev); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if ev.Type != "step_status" || ev.Step != "install-clo" {
		t.Fatalf("unexpected event: %+v", ev)
	}
}
