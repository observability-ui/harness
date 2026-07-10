package runner

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrefixWriter(t *testing.T) {
	var buf bytes.Buffer
	pw := newPrefixWriter(&buf, "backend", 0)

	pw.Write([]byte("starting server\n"))
	pw.Write([]byte("listening on :8080\n"))

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %q", len(lines), output)
	}
	for _, line := range lines {
		if !strings.Contains(line, "backend") {
			t.Fatalf("line should contain prefix 'backend': %q", line)
		}
	}
}
