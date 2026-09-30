package main

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

type FileData struct {
	AbsPath string
	Info    os.FileInfo
	Mime    string
	Data    []byte
}

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

func ValidateFile(path string) (file FileData, err error) {
	absPath, fileInfo, err := FindFile(path)
	if err != nil {
		return FileData{}, fmt.Errorf("invalid file: %w", err)
	}

	file.AbsPath = absPath
	file.Info = fileInfo

	data, err := os.ReadFile(file.AbsPath)
	if err != nil {
		return FileData{}, fmt.Errorf("unable to read file: %w", err)
	}

	file.Data = data

	nameParts := strings.Split(file.Info.Name(), ".")
	extension := "." + nameParts[len(nameParts)-1]
	mimeType := mime.TypeByExtension(extension)

	file.Mime = mimeType

	return
}
