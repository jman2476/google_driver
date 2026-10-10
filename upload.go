package main

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

func writeUploadMetadata(fileData FileData, targetFolderID string) *drive.File {
	meta := &drive.File{
		Name:     fileData.Info.Name(),
		MimeType: fileData.Mime,
	}

	if targetFolderID != "" {
		meta.Parents = []string{targetFolderID}
	}

	return meta
}

func (cfg *apiConfig) uploadFile(path string, targetDir string) (*drive.File, error) {
	fileData, err := ValidateFile(path)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to validate file: %v\n",
			err,
		)
	}
	targetId := "root"

	target, err := cfg.findFolder(targetDir)
	if err != nil && err != ErrEmptyDir {
		return nil, fmt.Errorf(
			"failed to find target folder: %w\n", err,
		)
	}
	if target != nil {
		targetId = target.Id
	}

	metadata := writeUploadMetadata(fileData, targetId)
	uploadBuilder := cfg.service.Files.Create(metadata)
	uploadBuilder.Media(fileData.Reader)

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

	fmt.Printf("File uploaded to %s\n", response.WebViewLink)

	return response, nil
}
