package process

import (
	"strings"
	"sync"
)

type ringBuffer struct {
	mu      sync.Mutex
	lines   []string
	maxSize int
	start   int
	count   int
	total   int
	partial string
}

func newRingBuffer(maxLines int) *ringBuffer {
	return &ringBuffer{
		lines:   make([]string, maxLines),
		maxSize: maxLines,
	}
}

func (rb *ringBuffer) Write(p []byte) (int, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	text := rb.partial + string(p)
	parts := strings.Split(text, "\n")

	// Last element is either empty (if text ended with \n) or a partial line
	rb.partial = parts[len(parts)-1]
	completedLines := parts[:len(parts)-1]

	for _, line := range completedLines {
		idx := (rb.start + rb.count) % rb.maxSize
		rb.lines[idx] = line
		if rb.count < rb.maxSize {
			rb.count++
		} else {
			rb.start = (rb.start + 1) % rb.maxSize
		}
		rb.total++
	}
	return len(p), nil
}

func (rb *ringBuffer) Lines() []string {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	result := make([]string, rb.count)
	for i := 0; i < rb.count; i++ {
		result[i] = rb.lines[(rb.start+i)%rb.maxSize]
	}
	if rb.partial != "" {
		result = append(result, rb.partial)
	}
	return result
}

func (rb *ringBuffer) Len() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return rb.total
}
