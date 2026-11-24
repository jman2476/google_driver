package main

import (
	"errors"
	"fmt"

	"google.golang.org/api/drive/v3"
)

var (
	ErrNeedRecursive       = errors.New("to download a folder, use the '-r' flag")
	ErrFileNotFound        = errors.New("no file found at this address")
	ErrDuplicateFilesFound = errors.New("multiple files found w/ same name")
)

func (cfg *apiConfig) downloadFile(path string, targetPath string, flags map[string]string) (*drive.File, error) {
	// look for file at drive location
	dirId, filename, isFolder, err := cfg.ValidateDrivePath(path)
	if err != nil {
		return nil, err
	}

	if isFolder {
		if _, ok := flags["r"]; ok {
			return cfg.downloadFolder(path, targetPath, flags)
		}
		return nil, ErrNeedRecursive
	}
	// if exists, download file
	// query drive for file
	// 		- get Id, get metadata
	query := fmt.Sprintf(
		"name = '%s' and '%s' in parents and trashed = false",
		filename,
		dirId,
	)

	files, err := cfg.service.Files.List().
		Q(query).
		Fields("files(id, name)").
		Do()

	if err != nil {
		return nil, err
	}

	if len(files.Files) == 0 {
		return nil, ErrFileNotFound
	}
	if len(files.Files) > 1 {
		return nil, ErrDuplicateFilesFound
	}

	fileId := files.Files[0].Id
	// use service.Files.Get(fileId).Download() to get data
	res, err := cfg.service.Files.Get(fileId).Download()
	if err != nil {
		return nil, fmt.Errorf(
			"Error downloading file: Status %d; %s\n",
			res.StatusCode, err,
		)
	}

	// check for targetPath location
	// if targetPath doesn't exist
	//		- write file at current directory
	// else
	//		- write file at target directory

	fmt.Println(dirId, filename)

	return nil, nil
}

func (cfg *apiConfig) downloadFolder(path string, targetPath string, flags map[string]string) (*drive.File, error) {

	return nil, nil
}
