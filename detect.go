package main

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"slices"
)

var (
	ErrFolderDNE = errors.New("no folder found at path")
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

// Function to read drive upload path and deliver:
// -- Target folder
// -- Target name
// -- is target a folder?

// What to check:
// 1- is the path just root?
// 2- does the path lead directly to a folder on drive?
// 3- does path end in '/', and we should only look for a folder?
// 4-
// 5-
// 6-
func (cfg apiConfig) ValidateUploadPath(path string) (
	dirID,
	name string,
	isFolder bool,
	err error,
) {
	roots := []string{
		"",
		"root",
		"/root",
		"root/",
		"/root/",
	}
	if slices.Contains(roots, path) {
		dirID = "root"
		isFolder = true
		return
	}

	folder, err := cfg.findFolder(path)
	if err == nil {
		dirID = folder.Id
		isFolder = true
		return
	} else if length := len(path); path[length-1] == '/' {
		err = fmt.Errorf("%w: %s", ErrFolderDNE, path)
		isFolder = true
		return
	}

	pathParts, err := ParsePath(path)
	if err != nil {
		return
	}

	idx := len(pathParts) - 1
	name = pathParts[idx].Name
	if len(pathParts) > 1 {
		dirPath := JoinFolderNames(pathParts[:idx])
		folder, err = cfg.findFolder(dirPath)
		if err != nil {
			err = fmt.Errorf("%w: %s", ErrFolderDNE, dirPath)
			isFolder = false
			return
		}
		dirID = folder.Id
	} else {
		dirID = "root"
	}

	return
}
