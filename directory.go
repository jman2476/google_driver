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

func (cfg *apiConfig) findFolder(path string) (*drive.File, error) {
	// Start assuming only search depth of one
	folders, err := ParsePath(path)
	if err != nil {
		if err == ErrEmptyDir {
			return nil, err
		}
		return nil, fmt.Errorf(
			"Find folder: path error: %v", err,
		)
	}

	err = cfg.validateToken()
	if err != nil {
		return nil, fmt.Errorf(
			"unable to validate token: %w", err,
		)
	}

	currentParentID := "root"
	var currentFolder *drive.File

	for _, folder := range folders {
		query := fmt.Sprintf(
			"mimeType = '%s' and name = '%s' and '%s' in parents and trashed = false",
			driveFolderMIME,
			folder.Name,
			currentParentID,
		)

		res, err := cfg.service.Files.List().
			Q(query).
			Fields("files(id, name, parents)").
			PageSize(1).
			Do()

		if err != nil {
			return nil, fmt.Errorf(
				"error querying folder '%s': %w",
				folder.Name, err,
			)
		}

		if len(res.Files) == 0 {
			return nil, fmt.Errorf(
				"%w: %s", ErrDirNotFound, folder.Name,
			)
		}

		currentFolder = res.Files[0]
		currentParentID = currentFolder.Id
	}

	return currentFolder, nil
}

func (cfg *apiConfig) createFolder(path string, _ string) (*drive.File, error) {
	folders, err := ParsePath(path)
	if err != nil && err != ErrEmptyDir {
		return nil, fmt.Errorf(
			"Create folder: path error: %w", err,
		)
	}

	parentId := "root"

	if len(folders) > 1 {
		var pathSlice []string
		for i, f := range folders {
			if i == len(folders)-1 {
				break
			}
			pathSlice = append(pathSlice, f.Name)
		}

		parentPath := strings.Join(pathSlice, "/")

		parent, err := cfg.findFolder(parentPath)
		if err != nil {
			return nil, fmt.Errorf(
				"parent of folder to create not found at %s: %w",
				parentPath, err,
			)
		}

		parentId = parent.Id
	}

	metadata := &drive.File{
		Name:     folders[len(folders)-1].Name,
		MimeType: driveFolderMIME,
		Parents:  []string{parentId},
	}
	uploadBuilder := cfg.service.Files.Create(metadata)

	err = cfg.validateToken()
	if err != nil {
		return nil, fmt.Errorf(
			"create folder: token validation error: %w", err,
		)
	}

	folderResp, err := uploadBuilder.Fields("id", "name", "webViewLink").Do()
	if err != nil {
		fmt.Printf("Response failure: %v\n", folderResp)

		return nil, fmt.Errorf(
			"create file error: %w", err,
		)
	}

	fmt.Printf("Folder made at %s\n", folderResp.WebViewLink)
	return folderResp, nil
}

func DeleteFolder(path string) (*drive.File, error) {
	return nil, nil
}
