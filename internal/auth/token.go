package auth

import (
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
