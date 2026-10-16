package main

import (
	"errors"
	"fmt"

	"google.golang.org/api/drive/v3"
)

var (
	ErrFileNotFound      = errors.New("File delete: file not found")
	ErrTooManyFilesFound = errors.New(
		"File delete: more than one file found with same name",
	)
)

func (cfg *apiConfig) removeFile(fileName, targetDir string) (*drive.File, error) {
	location, err := cfg.findFolder(targetDir)
	if err != nil {
		return nil, fmt.Errorf("unable to locate target directory: %w", err)
	}

	query := fmt.Sprintf(
		"name = '%s' and '%s' in parents and trashed = false",
		fileName, location.Id,
	)

	res, err := cfg.service.Files.List().
		Q(query).
		Fields("files(id, name)").
		Do()

	if err != nil {
		return nil, fmt.Errorf(
			"error finding file '%s': %w",
			fileName, err,
		)
	}

	if len(res.Files) == 0 {
		return nil, fmt.Errorf(
			"%w: %s", ErrFileNotFound, fileName,
		)
	}

	if len(res.Files) > 1 {
		return nil, fmt.Errorf(
			"%w: %s", ErrTooManyFilesFound, fileName,
		)
	}

	err = cfg.service.Files.Delete(res.Files[0].Id).Do()
	if err != nil {
		return nil, fmt.Errorf(
			"unable to delete file: %w", err,
		)
	}
	return nil, nil
}
