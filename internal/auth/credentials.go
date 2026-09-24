package auth

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func LoadCredentials() (*oauth2.Config, error) {
	path, err := filepath.Abs("./credentials.json")
	if err != nil {
		errMsg := fmt.Errorf(
			"could not get credential file path: %v", err,
		)
		return nil, errMsg
	}

	data, err := os.ReadFile(path)
	if err != nil {
		errMsg := fmt.Errorf(
		"error reading credentials.json: %w", err,
		)
		return nil, errMsg
	}

	// TODO: add option to choose between drive.DriveScope
	// 		and drive.DriveFileScope
	config, err := google.ConfigFromJSON(
		data, drive.DriveScope)
	if err != nil {
		errMsg := fmt.Errorf(
			"error parsing credentials: %w", err,
		)
		return nil, errMsg
	}

	return config, nil
}
