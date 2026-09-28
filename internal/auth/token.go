package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

const (
	authTokenPath      string = ".config/g-driver-portal/token.json"
	authCredentialPath string = ".config/g-driver-portal/credentials.json"
)

func ReadTokenData() (*oauth2.Token, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf(
			"unable to find user home directory: %w", err,
		)
	}
	absPath := filepath.Join(homeDir, authTokenPath)
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to read auth token: %w", err,
		)
	}

	var token *oauth2.Token
	err = json.Unmarshal(data, &token)
	if err != nil {
		return nil, fmt.Errorf(
			"error reading token: %w", err,
		)
	}

	return token, nil
}

func WriteTokenData(token *oauth2.Token) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf(
			"unable to find user home directory: %w", err,
		)
	}
	dirPath := filepath.Join(homeDir, authTokenPath)
	fmt.Printf("Directory path: %s\n", dirPath)
	err = os.MkdirAll(dirPath, 0700)
	if err != nil {
		return fmt.Errorf("unable to make config folder: %w", err)
	}

	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("token marshalling error: %w", err)
	}

	filePath := dirPath + "/token.json"
	fmt.Printf("File path: %s\n", filePath)
	err = os.WriteFile(filePath, data, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache token data: %w", err)
	}

	return nil
}

func GetAuthorization(c *oauth2.Config) (authCode string, err error) {
	authURL := c.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	)

	fmt.Println("To log in, copy the URL below and paste it into your browser:")
	fmt.Println(authURL)

	fmt.Println("Getting your authorization code...")
	authCode, err = GetAuthCode()
	if err != nil {
		return "", fmt.Errorf("unable to read authorization code from server: %w", err)
	}

	return
}

func GetToken(c *oauth2.Config, code string) (token *oauth2.Token, err error) {
	token, err = c.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("unable to exchange authorization code for access token: %w", err)
	}

	err = WriteTokenData(token)
	if err != nil {
		return token, fmt.Errorf("Error writing token data: %w", err)
	}

	return
}
