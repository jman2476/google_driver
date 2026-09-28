package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jman2476/google-driver/internal/auth"
)

func (cfg *apiConfig) setClient() error {
	log.Println("Starting Google Driver")
	config, err := auth.LoadCredentials()
	if err != nil {
		return fmt.Errorf("Error loading credentials: %w\n", err)
	}

	token, err := auth.ReadTokenData()
	if err != nil {
		fmt.Println("No token found. Initiating authorization.")

		authCode, err := auth.GetAuthorization(config)
		if err != nil {
			return fmt.Errorf("unable to authorize: %w", err)
		}

		token, err = auth.GetToken(config, authCode)
		if err != nil {
			return fmt.Errorf("Error getting token: %w", err)
		}
	}

	cfg.client = config.Client(context.Background(), token)
	cfg.token = token

	return nil
}
