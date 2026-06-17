package state

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type Lock struct {
	path string
	file *os.File
}

func NewLock(stateDir string) *Lock {
	return &Lock{path: filepath.Join(stateDir, "obsui.lock")}
}

func (l *Lock) Acquire() error {
	os.MkdirAll(filepath.Dir(l.path), 0755)
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("cannot open lock file: %w", err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		return fmt.Errorf("another obsui instance is already running (lock: %s)", l.path)
	}
	l.file = f
	return nil
}

func (l *Lock) Release() {
	if l.file != nil {
		syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
		l.file.Close()
		os.Remove(l.path)
	}
}
