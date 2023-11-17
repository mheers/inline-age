package helpers

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func Root() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get caller information")
	}

	dir := filepath.Dir(filename)
	file := filepath.Join(dir, "..")
	absPath, err := filepath.Abs(file)
	if err != nil {
		return "", err

	}
	return absPath, nil
}
