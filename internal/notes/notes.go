package notes

import (
	"fmt"
	"path/filepath"

	"github.com/rafawastaken/quicktick/internal/config"
)

func GetNotePath(id int64) (string, error) {
	dir, err := config.NotesDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, fmt.Sprintf("%d.md", id)), nil
}
