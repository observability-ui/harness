package output

import (
	"encoding/json"
	"io"
	"sync"
)

type Event struct {
	Type   string `json:"type"`
	Step   string `json:"step,omitempty"`
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type JSONEmitter struct {
	mu  sync.Mutex
	enc *json.Encoder
}

func NewJSONEmitter(w io.Writer) *JSONEmitter {
	return &JSONEmitter{enc: json.NewEncoder(w)}
}

func (e *JSONEmitter) Emit(ev Event) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.enc.Encode(ev)
}
