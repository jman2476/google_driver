package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

var (
	errTokenExpired    = errors.New("token has expired")
	errTokenNearExpiry = errors.New("token expires soon")
)

func NeedRefresh(token *oauth2.Token, buffer time.Duration) bool {

	if token == nil || !token.Valid() {
		return true
	}

	return time.Until(token.Expiry) <= buffer
}

func RefreshToken(config *oauth2.Config, token *oauth2.Token) (*oauth2.Token, error) {
	tokenSource := config.TokenSource(context.Background(), token)

	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return newToken, nil
}

func EnsureValidToken(config *oauth2.Config, token *oauth2.Token) (*oauth2.Token, bool, error) {
	if !NeedRefresh(token, 5*time.Minute) {
		return token, false, nil
	}

	newToken, err := RefreshToken(config, token)
	if err != nil {
		return nil, false, err
	}

	wasRefreshed := newToken.AccessToken != token.AccessToken

	return newToken, wasRefreshed, nil
}
