package auth

import (
	"fmt"
	"path/filepath"

	"golang.org/x/oauth2"
)

func LoadCredentials() (oauth2.Config, error) {
	path, err := filepath.Abs("./credentials.json")
	if err != nil {
		errMsg := fmt.Errorf("Could not get credential file: %v", err)
		return oauth2.Config{}, errMsg
	}

}
