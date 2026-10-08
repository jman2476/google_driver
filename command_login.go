package main

import (
	"context"
	"fmt"

	"github.com/jman2476/google-driver/internal/auth"
	"google.golang.org/api/drive/v3"
)

func (cfg *apiConfig) commandLogIn(_, _ string) (*drive.File, error) {
	config, err := auth.LoadCredentials()
	if err != nil {
		return nil, fmt.Errorf("unable to load credentials: %w\n", err)
	}

	authCode, err := auth.GetAuthorization(config)
	if err != nil {
		return nil, fmt.Errorf("unable to authorize: %w\n", err)
	}

	token, err := auth.GetToken(config, authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to get token: %w\n", err)
	}

	cfg.client = config.Client(context.Background(), token)
	cfg.token = token

	err = cfg.setService()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (cfg *apiConfig) commandLogOut(_, _ string) (*drive.File, error) {
	return nil, cfg.eraseTokenCache()
}
