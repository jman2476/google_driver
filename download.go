package main

import (
	"errors"
	"fmt"

	"google.golang.org/api/drive/v3"
)

var (
	ErrNeedRecursive = errors.New("to download a folder, use the '-r' flag")
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
	// use service.Files.Get(fileId).Download() to get data
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
