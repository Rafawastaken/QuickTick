package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rafawastaken/quicktick/internal/config"
)

type Session struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	Email       string `json:"email"`
}

func AuthPath() (string, error) {
	dir, err := config.AppDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "auth.json"), nil
}

func SaveSession(s *Session) error {
	path, err := AuthPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600) // Secure permissions
}

func LoadSession() (*Session, error) {
	path, err := AuthPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil // No session
	}
	if err != nil {
		return nil, err
	}

	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func ClearSession() error {
	path, err := AuthPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
