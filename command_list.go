package main

import (
	"fmt"

	"github.com/jman2476/google-driver/internal/progress"
	"google.golang.org/api/drive/v3"
)

func (cfg *apiConfig) commandList(targetDir, _ string) (*drive.File, error) {
	target, err := cfg.findFolder(targetDir)
	if err != nil {
		return nil, fmt.Errorf("Failed to find folder: %w", err)
	}

	contents, err := cfg.listFolder(target.Id)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to list folder contents: %w", err,
		)
	}

	if contents.Files == nil {
		return nil, nil
	}
	printFiles(contents.Files)

	return target, nil
}

func printFiles(fileList []*drive.File) {
	if fileList == nil {
		return
	}
	for _, f := range fileList {
		fmt.Printf("%v -- %v -- %v\n", f.Name, f.MimeType, progress.FormatSize(f.Size))
	}
}
