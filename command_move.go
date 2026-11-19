package main

import (
	"google.golang.org/api/drive/v3"
)

func (cfg *apiConfig) commandMove(args arguments) (*drive.File, error) {
	if len(args.positional) < 2 {
		return nil, ErrMissingArgs
	}
	if len(args.positional) > 2 {
		return nil, ErrTooManyArgs
	}

	source := args.positional[0]
	target := args.positional[1]

	return cfg.uploadFile(source, target)
}
