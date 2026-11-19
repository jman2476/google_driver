package main

import "google.golang.org/api/drive/v3"

func (cfg *apiConfig) commandDownload(args arguments) (*drive.File, error) {
	var source, target string
	switch len(args.positional) {
	case 0:
	case 1:
		return nil, ErrMissingArgs
	case 2:

	}

	return cfg.downloadFile(source, target)
}
