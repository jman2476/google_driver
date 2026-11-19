package main

import "google.golang.org/api/drive/v3"

func (cfg *apiConfig) commandCreateFolder(args arguments) (*drive.File, error) {
	switch len(args.positional) {
	case 0:
		return nil, ErrMissingArgs
	case 1:
		return cfg.createFolder(args.positional[0])
	default:
		return nil, ErrTooManyArgs
	}
}
