package main

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/api/drive/v3"
)

const (
	driveFolderMIME = "application/vnd.google-apps.folder"
)

var (
	ErrEmptyDir    = errors.New("Directory path: empty path")
	ErrParseDirs   = errors.New("Directory path: unable to parse")
	ErrDirNotFound = errors.New("Directory path: not found on Drive")
)

type DriveFolder struct {
	Name   string
	Parent string
}

func ParsePath(folderPath string) ([]DriveFolder, error) {
	folders := strings.Split(folderPath, "/")

	var parsedFolders []DriveFolder
	for i, f := range folders {
		if f != "" {
			df := DriveFolder{Name: f}
			if i > 0 {
				df.Parent = folders[i-1]
			} else {
				df.Parent = ""
			}
			parsedFolders = append(parsedFolders, df)
		} else if i > 0 && i < len(folders)-1 {
			return []DriveFolder{}, ErrParseDirs
		}
	}

	if len(parsedFolders) == 0 {
		return []DriveFolder{}, ErrEmptyDir
	}

	return parsedFolders, nil
}

func FindFolder(path string) (*drive.File, error) {
	// Start assuming only search depth of one
	folders := strings.Split(path, "/")

	fmt.Printf("folders: %s", folders)
	searchQuery := fmt.Sprintf("mimeType = '%s'", driveFolderMIME)
	fmt.Printf("Drive MIME type: %s", searchQuery)

	return nil, nil
}

func CreateFolder(path string) (*drive.File, error) {
	return nil, nil
}

func DeleteFolder(path string) (*drive.File, error) {
	return nil, nil
}
