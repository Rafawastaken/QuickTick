package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const AppDirName = "quicktick"

func AppDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, AppDirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func DBPath(userID string) (string, error) {
	dir, err := AppDir()
	if err != nil {
		return "", err
	}
	// If no user ID (e.g. guest or pre-login), use default
	filename := "quicktick.db"
	if userID != "" {
		filename = fmt.Sprintf("quicktick_%s.db", userID)
	}
	return filepath.Join(dir, filename), nil
}

func NotesDir() (string, error) {
	dir, err := AppDir()
	if err != nil {
		return "", err
	}
	nd := filepath.Join(dir, "notes")
	if err := os.MkdirAll(nd, 0o755); err != nil {
		return "", err
	}
	return nd, nil
}
