package state

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// Lock is released implicitly when the process exits and the file descriptor closes.

type Lock struct {
	path string
	file *os.File
}

func NewLock(stateDir string) *Lock {
	return &Lock{path: filepath.Join(stateDir, "obs.lock")}
}

func (l *Lock) Acquire() error {
	if err := os.MkdirAll(filepath.Dir(l.path), 0755); err != nil {
		return fmt.Errorf("cannot create lock directory: %w", err)
	}
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("cannot open lock file: %w", err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		return fmt.Errorf("another obs instance is already running (lock: %s)", l.path)
	}
	l.file = f
	return nil
}

func (l *Lock) Release() error {
	if l.file == nil {
		return nil
	}
	syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	err := l.file.Close()
	l.file = nil
	os.Remove(l.path)
	return err
}

