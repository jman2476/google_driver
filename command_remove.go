package main

import (
	"strings"

	"google.golang.org/api/drive/v3"
)

func (cfg *apiConfig) commandRemove(args arguments) (*drive.File, error) {
	var filepath string

	switch len(args.positional) {
	case 0:
		return nil, ErrMissingArgs
	case 1:
		filepath = args.positional[0]
	default:
		return nil, ErrTooManyArgs
	}

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
