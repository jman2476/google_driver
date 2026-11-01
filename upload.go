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

	fmt.Printf("File size: %s\n\n", progress.FormatSize(fileData.Info.Size()))

	progressTracker := &progress.Tracker{
		File: fileData.Reader}
	progressTracker.SetSize(fileData.Info.Size())
	progressTracker.SetTime(time.Now())
	metadata := writeUploadMetadata(fileData, targetId, targetName)
	uploadBuilder := cfg.service.Files.Create(metadata)
	uploadBuilder.Media(progressTracker)

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
