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

	// TODO: add option to choose between
	// 		drive.DriveScope -> full access
	// 		and
	// 		drive.DriveFileScope -> restricted access
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

func PrintConfig(c *oauth2.Config) {
	fmt.Printf("Client ID: %v\n", c.ClientID)
	fmt.Printf("Client Secret: %v\n", c.ClientSecret)
	fmt.Printf("Client Endpoint: %v\n", c.Endpoint)
	fmt.Printf("Redirect URL: %v\n", c.RedirectURL)
	fmt.Printf("Scopes: %v\n", c.Scopes)
}
