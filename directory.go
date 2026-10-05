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

func printDriveFolder(folders []DriveFolder) {
	for _, f := range folders {
		fmt.Printf("Name: %v, Parent: %v\n", f.Name, f.Parent)
	}
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

func (cfg *apiConfig) FindFolder(path string) (*drive.File, error) {
	// Start assuming only search depth of one
	folders, err := ParsePath(path)
	if err != nil {
		return nil, fmt.Errorf(
			"Find folder: path error: %v", err,
		)
	}

	printDriveFolder(folders)

	fmt.Printf("folders: %v", folders)
	searchQuery := fmt.Sprintf("mimeType = '%s'", driveFolderMIME)
	fmt.Printf("Drive MIME type: %s", searchQuery)

	folderList, err := cfg.service.Files.List().Q(searchQuery).Do()
	if err != nil {
		return nil, fmt.Errorf(
			"FindFolder error: %w", err,
		)
	}

	for _, f := range folderList.Files {
		fmt.Printf("Folder name: %v ID: %v  Parent: %v", f.Name, f.Id, f.Parents[0])
	}

	var tracker = struct {
		ParentID    string
		ParentName  string
		CurrentFile *drive.File
	}{}
	for _, folder := range folders {
		fmt.Printf("Checking for folder: %s\n", folder.Name)
		for _, f := range folderList.Files {
			fmt.Printf("Folder: %s\n", f.Name)
			if folder.Name == f.Name {
				if folder.Parent == "" {
					tracker.CurrentFile = f
					break
				} else if folder.Parent == f.Parents[0] {
					tracker.ParentName = tracker.CurrentFile.Name
					tracker.CurrentFile = f
					tracker.ParentID = f.Parents[0]
					break
				}
			}
		}
	}

	if tracker.CurrentFile == nil {
		return nil, fmt.Errorf("No folder found")
	}

	return tracker.CurrentFile, nil
}

func CreateFolder(path string) (*drive.File, error) {
	return nil, nil
}

func DeleteFolder(path string) (*drive.File, error) {
	return nil, nil
}
