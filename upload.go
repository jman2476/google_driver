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

func UploadFile(fs *drive.FilesService, path string, targetDir string) error {
	fileData, err := ValidateFile(path)
	if err != nil {
		return fmt.Errorf(
			"failed to validate file: %v",
			err,
		)
	}

	metadata := WriteUploadMetadata(fileData, targetDir)
	uploadBuilder := fs.Create(metadata)
	uploadBuilder.Media(fileData.Reader)

	return nil
}
