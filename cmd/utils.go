package cmd

import (
	"fmt"
	"os"
	"strings"

	Errors "github.com/pkg/errors"
)

func isSymlink(path string) (bool, error) {
	file, err := os.Lstat(path)
	if Errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Creating symlink from %s\n", path)
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return file.Mode()&os.ModeSymlink != 0, nil
}

func removeHomeDir(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = strings.Replace(path, "~", homeDir, 1)
	}
	return path, nil
}

func pathType(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "dne", nil
		}
		return "", err
	}
	if info.IsDir() {
		return "directory", nil
	}
	return "file", nil
}
