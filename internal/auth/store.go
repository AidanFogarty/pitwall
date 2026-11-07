package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ErrNotLoggedIn = errors.New("not logged in")

type Store struct {
	path string
}

type storedSession struct {
	Token string `json:"subscription_token"`
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

func DefaultStore() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("finding home directory: %w", err)
	}
	return NewStore(filepath.Join(home, ".pitwall", "auth.json")), nil
}

func (s *Store) Save(session *Session) error {
	data, err := json.Marshal(storedSession{Token: session.Token})
	if err != nil {
		return fmt.Errorf("encoding session: %w", err)
	}

	if err := writeFileAtomic(s.path, data, 0o600); err != nil {
		return fmt.Errorf("saving session to %s: %w", s.path, err)
	}
	return nil
}

func (s *Store) Load() (*Session, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotLoggedIn
		}
		return nil, fmt.Errorf("reading %s: %w", s.path, err)
	}

	var stored storedSession
	if err := json.Unmarshal(data, &stored); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", s.path, err)
	}

	return sessionFromToken(stored.Token)
}

func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating %s: %w", dir, err)
	}

	tmp, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	if err := writeAndClose(tmp, data, perm); err != nil {
		return err
	}

	return os.Rename(tmp.Name(), path)
}

func writeAndClose(file *os.File, data []byte, perm os.FileMode) (err error) {
	defer func() { err = errors.Join(err, file.Close()) }()

	if err := file.Chmod(perm); err != nil {
		return fmt.Errorf("setting permissions: %w", err)
	}
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("writing: %w", err)
	}
	return nil
}
