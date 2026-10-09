package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

type FileData struct {
	AbsPath string
	Info    os.FileInfo
	Mime    string
	Reader  *os.File
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

	fileHandle, err := os.Open(file.AbsPath)
	if err != nil {
		return FileData{}, fmt.Errorf("unable to read file: %w", err)
	}

	file.Reader = fileHandle

	// nameParts := strings.Split(file.Info.Name(), ".")
	// extension := "." + nameParts[len(nameParts)-1]
	// mimeType := mime.TypeByExtension(extension)

	extension := filepath.Ext(file.Info.Name())
	mimeType := mime.TypeByExtension(extension)
	if mimeType == "" {
		buf := make([]byte, 512)
		n, _ := fileHandle.Read(buf)
		fileHandle.Seek(0, io.SeekStart)
		if n > 0 {
			mimeType = http.DetectContentType(buf[:n])
		} else {
			mimeType = "application/octet-stream"
		}
	}

	file.Mime = mimeType
	return
}
