package auth

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

var (
	errTokenExpired    = errors.New("token has expired")
	errTokenNearExpiry = errors.New("token expires soon")
)

func checkTokenAge(t *oauth2.Token) error {

	if !t.Expiry.After(time.Now()) {
		return errTokenExpired
	}

	dayDuration, err := time.ParseDuration("24h")
	if err != nil {
		return fmt.Errorf("duration parse error: %v", err)
	}
	if time.Until(t.Expiry) <= dayDuration {
		return errTokenNearExpiry
	}

	return nil
}
