package main

import (
	"strings"

	"google.golang.org/api/drive/v3"
)

func (cfg *apiConfig) commandRemove(filepath, _ string) (*drive.File, error) {
	lastSlashIdx := strings.LastIndex(filepath, "/")
	var folderpath string
	var filename string
	switch lastSlashIdx {
	case -1:
		folderpath = "root"
		filename = filepath
	default:
		folderpath = filepath[:lastSlashIdx]
		filename = filepath[lastSlashIdx+1:]
	}

	return cfg.removeFile(filename, folderpath)
}
