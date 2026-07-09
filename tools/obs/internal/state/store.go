package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"obs/internal/process"
)

type ProcessState struct {
	Name string `json:"name"`
	PID  int    `json:"pid"`
}

type RunState struct {
	Processes []ProcessState `json:"processes"`
}

type Store struct {
	dir         string
	mu          sync.Mutex
	initialized bool
}

func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

func (s *Store) init() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.initialized {
		return nil
	}
	for _, sub := range []string{"", "pids"} {
		if err := os.MkdirAll(filepath.Join(s.dir, sub), 0755); err != nil {
			return err
		}
	}
	s.initialized = true
	return nil
}

func (s *Store) Save(state *RunState) error {
	if err := s.init(); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	tmp := filepath.Join(s.dir, "state.json.tmp")
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(s.dir, "state.json"))
}

func (s *Store) Load() (*RunState, error) {
	data, err := os.ReadFile(filepath.Join(s.dir, "state.json"))
	if err != nil {
		return nil, err
	}
	var state RunState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (s *Store) Clean() {
	os.Remove(filepath.Join(s.dir, "state.json"))
	os.Remove(filepath.Join(s.dir, "state.json.tmp"))
	os.Remove(filepath.Join(s.dir, "obs.lock"))
	os.RemoveAll(filepath.Join(s.dir, "pids"))
}

func (s *Store) WritePID(name string, pid int) error {
	if err := s.init(); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(s.dir, "pids", name+".pid"), []byte(strconv.Itoa(pid)), 0644)
}

func DefaultStateDir() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "tools")); err == nil {
			return filepath.Join(dir, ".obs")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, ".obs")
}

func (s *Store) FilterAlive(rs *RunState) *RunState {
	var alive []ProcessState
	for _, p := range rs.Processes {
		if process.IsAlive(p.PID) {
			alive = append(alive, p)
		}
	}
	return &RunState{Processes: alive}
}

func PrintStatus(rs *RunState) {
	if len(rs.Processes) == 0 {
		fmt.Println("No running processes.")
		return
	}
	fmt.Printf("%-25s %-8s %s\n", "PROCESS", "PID", "STATUS")
	for _, p := range rs.Processes {
		alive := "dead"
		if process.IsAlive(p.PID) {
			alive = "running"
		}
		fmt.Printf("%-25s %-8d %s\n", p.Name, p.PID, alive)
	}
}
