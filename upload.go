package main

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

func WriteUploadMetadata(fileData FileData, targetFolderID string) *drive.File {
	meta := &drive.File{
		Name:     fileData.Info.Name(),
		MimeType: fileData.Mime,
	}

	if targetFolderID != "" {
		meta.Parents = []string{targetFolderID}
	}

	return meta
}

func (cfg *apiConfig) UploadFile(path string, targetDir string) error {
	fileData, err := ValidateFile(path)
	if err != nil {
		return fmt.Errorf(
			"failed to validate file: %v",
			err,
		)
	}

	metadata := WriteUploadMetadata(fileData, targetDir)
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
		return fmt.Errorf(
			"upload file: token validation error: %w", err,
		)
	}

	response, err := uploadBuilder.Fields("id", "name").Do()
	if err != nil {
		fmt.Printf("Response failure: %v", response)

		return fmt.Errorf(
			"upload file error: %w", err,
		)
	}

	fmt.Printf("Response success: %v", response)

	return nil
}
