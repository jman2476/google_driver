package main

import (
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
