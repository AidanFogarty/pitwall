package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Logger struct {
	mu sync.Mutex
	w  io.Writer
}

func NewDiscard() *Logger {
	return &Logger{w: io.Discard}
}

func NewFile(path string) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	return &Logger{w: f}, nil
}

func DefaultPath(name string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return name
	}

	return filepath.Join(home, ".pitwall", name)
}

func (l *Logger) Log(keyVals ...any) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	line := "ts=" + time.Now().UTC().Format(time.RFC3339Nano)
	for i := 0; i+1 < len(keyVals); i += 2 {
		line += fmt.Sprintf(" %v=%v", keyVals[i], keyVals[i+1])
	}

	_, err := fmt.Fprintln(l.w, line)
	return err
}

func (l *Logger) Logf(format string, args ...any) {
	_ = l.Log("src", "pitwall", "msg", fmt.Sprintf(format, args...))
}
