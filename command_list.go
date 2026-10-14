package main

import "google.golang.org/api/drive/v3"

func (cfg *apiConfig) commandList(_, targetDir string) (*drive.File, error) {
	target, err := cfg.findFolder()

	return nil, nil
}
