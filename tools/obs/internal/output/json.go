package output

import (
	"encoding/json"
	"io"
	"sync"
)

type Event struct {
	Type    string `json:"type"`
	Step    string `json:"step,omitempty"`
	Process string `json:"process,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	PID     int    `json:"pid,omitempty"`
}

type JSONEmitter struct {
	mu  sync.Mutex
	enc *json.Encoder
}

func NewJSONEmitter(w io.Writer) *JSONEmitter {
	return &JSONEmitter{enc: json.NewEncoder(w)}
}

func (e *JSONEmitter) Emit(ev Event) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.enc.Encode(ev)
}
