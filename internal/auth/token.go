package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
)

const (
	authTokenPath string = "~/.config/g-driver-portal/token.json"
)

func ReadTokenData() (*oauth2.Token, error) {
	data, err := os.ReadFile(authTokenPath)
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

func GetAuthorization(c *oauth2.Config) (authCode string, err error) {
	authURL := c.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	)

	fmt.Println("To log in, copy the URL below and paste it into your browser:")
	fmt.Println(authURL)

	fmt.Println("Enter your authorization code: ")
	_, err = fmt.Scan(&authCode)
	if err != nil {
		return "", fmt.Errorf("unable to read authorization code: %w", err)
	}

	return
}

func GetToken(c *oauth2.Config, code string) (token *oauth2.Token, err error) {
	token, err = c.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("unable to exchange authorization code for access token: %w", err)
	}

	return
}
