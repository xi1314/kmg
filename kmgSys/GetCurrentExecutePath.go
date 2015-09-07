package kmgSys

import (
	"path/filepath"
)

func GetCurrentExecutePath() (string, error) {
	p, err := getCurrentExecutePath()
	return filepath.Clean(p), err
}
