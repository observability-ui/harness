package process_test

import (
	"testing"

	"obs/internal/process"
)

func TestRingBuffer_Basic(t *testing.T) {
	rb := process.NewRingBuffer(3)
	rb.Write([]byte("line1\nline2\nline3\n"))

	lines := rb.Lines()
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(lines), lines)
	}
	if lines[0] != "line1" || lines[2] != "line3" {
		t.Fatalf("unexpected lines: %v", lines)
	}
}

func TestRingBuffer_Overflow(t *testing.T) {
	rb := process.NewRingBuffer(2)
	rb.Write([]byte("a\nb\nc\n"))

	lines := rb.Lines()
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "b" || lines[1] != "c" {
		t.Fatalf("oldest line should be evicted: %v", lines)
	}
}

func TestRingBuffer_PartialLine(t *testing.T) {
	rb := process.NewRingBuffer(5)
	rb.Write([]byte("hel"))
	rb.Write([]byte("lo\nworld\n"))

	lines := rb.Lines()
	if len(lines) != 2 || lines[0] != "hello" {
		t.Fatalf("partial writes should be joined: %v", lines)
	}
}
