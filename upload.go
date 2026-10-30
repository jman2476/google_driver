package main

import (
	"fmt"
	"time"

	"github.com/jman2476/google-driver/internal/progress"
	"google.golang.org/api/drive/v3"
)

func writeUploadMetadata(fileData FileData, targetFolderID, newName string) *drive.File {
	meta := &drive.File{
		Name:     fileData.Info.Name(),
		MimeType: fileData.Mime,
	}

	if targetFolderID != "" {
		meta.Parents = []string{targetFolderID}
	}

	if newName != "" {
		meta.Name = newName
	}

	return meta
}

func (cfg *apiConfig) uploadFile(path string, targetPath string) (*drive.File, error) {
	fileData, err := ValidateFile(path)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to validate file: %v\n",
			err,
		)
	}

	targetId, targetName, _, err := cfg.ValidateUploadPath(targetPath)
	if err != nil {
		return nil, fmt.Errorf(
			"bad target path: %v\n",
			err,
		)
	}
	defer fileData.Reader.Close()

	// targetId := "root"
	// if targetPath == "" {
	// 	targetPath = targetId
	// }

	// pathParts, err := ParsePath(targetPath)
	// if err != nil {
	// 	return nil, fmt.Errorf(
	// 		"failed to parse target path: %w\n", err,
	// 	)
	// }
	// var pathPartNames []string
	// for _, p := range pathParts {
	// 	pathPartNames = append(pathPartNames, p.Name)
	// }

	// var targetDir string
	// var targetName string
	// if lengthStr := len(targetPath); lengthStr > 0 && targetPath[lengthStr-1] == '/' {
	// 	targetDir = strings.Join(pathPartNames, "/")
	// 	targetName = fileData.Info.Name()
	// } else if length := len(pathParts); strings.Contains(pathPartNames[length-1], ".") {
	// 	targetDir = strings.Join(pathPartNames[:length-1], "/")
	// 	targetName = pathPartNames[length-1]
	// }

	// target, err := cfg.findFolder(targetDir)
	// if err != nil && err != ErrEmptyDir {
	// 	return nil, fmt.Errorf(
	// 		"failed to find target folder: %w\n", err,
	// 	)
	// }
	// if target != nil {
	// 	targetId = target.Id
	// }

	fmt.Printf("File size: %s\n\n", progress.FormatSize(fileData.Info.Size()))

	progressTracker := &progress.Tracker{
		File: fileData.Reader}
	progressTracker.SetSize(fileData.Info.Size())
	progressTracker.SetTime(time.Now())
	metadata := writeUploadMetadata(fileData, targetId, targetName)
	uploadBuilder := cfg.service.Files.Create(metadata)
	uploadBuilder.Media(progressTracker)

	// TODO:
	// - Check if target directory exists/directory path exists
	// - If it doesn't exist, create that directory
	// 		- directories on google drive are just files
	// 		  with MIME type application/vnd.google-apps.folder
	// - Maybe prompt user that the selected folder doesn't exist?
	// - if all is good, upload

	err = cfg.validateToken()
	if err != nil {
		return nil, fmt.Errorf(
			"upload file: token validation error: %w\n", err,
		)
	}

	response, err := uploadBuilder.Fields("id", "name", "parents", "size", "webViewLink").Do()
	if err != nil {
		fmt.Printf("Response failure: %v\n", response)

		return nil, fmt.Errorf(
			"upload file error: %w\n", err,
		)
	}

	fmt.Printf("\nFile uploaded to %s\n", response.WebViewLink)

	return response, nil
}
