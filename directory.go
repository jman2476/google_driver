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

	// printDriveFolder(folders)

	// fmt.Printf("folders: %v", folders)
	// searchQuery := fmt.Sprintf("mimeType = '%s'", driveFolderMIME)
	// fmt.Printf("Drive MIME type: %s", searchQuery)
	// var fields googleapi.Field = "files(id, name, parents)"

	// folderList, err := cfg.service.Files.List().Q(searchQuery).Fields(fields).Do()
	// if err != nil {
	// 	return nil, fmt.Errorf(
	// 		"FindFolder error: %w", err,
	// 	)
	// }

	// for _, f := range folderList.Files {
	// 	fmt.Printf("Folder name: %v ID: %v  Parent: %v", f.Name, f.Id, f.Parents[0])
	// }

	// var tracker = struct {
	// 	ParentID    string
	// 	ParentName  string
	// 	CurrentFile *drive.File
	// }{}
	// for _, folder := range folders {
	// 	fmt.Printf("Checking for folder: %s\n", folder.Name)
	// 	for _, f := range folderList.Files {
	// 		fmt.Printf("Folder: %s\n", f.Name)
	// 		if folder.Name == f.Name {
	// 			if folder.Parent == "" {
	// 				tracker.CurrentFile = f
	// 				break
	// 			} else if folder.Parent == f.Parents[0] {
	// 				tracker.ParentName = tracker.CurrentFile.Name
	// 				tracker.CurrentFile = f
	// 				tracker.ParentID = f.Parents[0]
	// 				break
	// 			}
	// 		}
	// 	}
	// }

	// if tracker.CurrentFile == nil {
	// 	return nil, ErrDirNotFound
	// }

	// return tracker.CurrentFile, nil
}

func (cfg *apiConfig) CreateFolder(path string) (*drive.File, error) {
	folders, err := ParsePath(path)
	if err != nil {
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

		parent, err := cfg.FindFolder(parentPath)
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

	folderResp, err := uploadBuilder.Fields("id", "name").Do()
	if err != nil {
		fmt.Printf("Response failure: %v", folderResp)

		return nil, fmt.Errorf(
			"create file error: %w", err,
		)
	}

	fmt.Printf("Response success: %v", folderResp)
	return folderResp, nil
}

func DeleteFolder(path string) (*drive.File, error) {
	return nil, nil
}
