package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindFile(path string) (string, os.FileInfo, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", nil, fmt.Errorf(
			"unable to construct absolute path for %v: %w",
			path, err,
		)
	}

	fileInfo, err := os.Lstat(absolutePath)
	if err != nil {
		return "", nil, fmt.Errorf(
			"cannot read file information: %w", err,
		)
	}

	return absolutePath, fileInfo, nil
}
