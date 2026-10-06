package auth

import (
	"errors"
	"time"

	"golang.org/x/oauth2"
)

var (
	errTokenExpired = errors.New("")
)

func checkTokenAge(t *oauth2.Token) error {

	if t.Expiry.Compare(time.Now()) >= 0 {
		return errTokenExpired
	}

	return nil
}
