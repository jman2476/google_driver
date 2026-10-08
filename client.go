package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jman2476/google-driver/internal/auth"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func (cfg *apiConfig) setClient() error {
	log.Println("Starting Google Driver")
	config, err := auth.LoadCredentials()
	if err != nil {
		return fmt.Errorf("Error loading credentials: %w\n", err)
	}

	token, err := auth.ReadTokenData()
	if err != nil {
		return fmt.Errorf(
			"no token found: %w\n please log in\n", err)
		// fmt.Println("No token found. Initiating authorization.")

		// authCode, err := auth.GetAuthorization(config)
		// if err != nil {
		// 	return fmt.Errorf("unable to authorize: %w", err)
		// }

		// token, err = auth.GetToken(config, authCode)
		// if err != nil {
		// 	return fmt.Errorf("Error getting token: %w", err)
		// }
	}

	cfg.client = config.Client(context.Background(), token)
	cfg.token = token

	return nil
}

func (cfg *apiConfig) setService() error {
	service, err := drive.NewService(
		context.Background(),
		option.WithHTTPClient(cfg.client),
	)
	if err != nil {
		return fmt.Errorf(
			"unable to start Google Drive service: %w", err,
		)
	}

	cfg.service = service

	return nil
}

func (cfg *apiConfig) validateToken() error {
	config, err := auth.LoadCredentials()
	if err != nil {
		return fmt.Errorf("error loading credentials: %w\n", err)
	}

	newToken, isRefresh, err := auth.EnsureValidToken(config, cfg.token)
	if err != nil {
		return fmt.Errorf("token validation error: %w\n", err)
	}

	if isRefresh {
		err = auth.WriteTokenData(newToken)
		if err != nil {
			return fmt.Errorf("unable to save new valid token: %w", err)
		}

		cfg.token = newToken
	}

	return nil
}

func (cfg *apiConfig) eraseTokenCache() error {
	return auth.DeleteTokenData()
}
