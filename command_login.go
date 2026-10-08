package main

import "google.golang.org/api/drive/v3"

func (cfg *apiConfig) commandLogIn(_, _ string) (*drive.File, error) {
	return nil, nil
}

func (cfg *apiConfig) commandLogOut(_, _ string) (*drive.File, error) {
	return nil, cfg.eraseTokenCache()
}
