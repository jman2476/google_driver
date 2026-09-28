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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf(
			"unable to find user home directory: %w", err,
		)
	}

	credentialPath := filepath.Join(homeDir, authCredentialPath)
	data, err := os.ReadFile(credentialPath)
	if err != nil {
		credentialPath = "./credentials.json"
		data, err = os.ReadFile(credentialPath)
	}

	if err != nil {
		return nil, fmt.Errorf(
			"no credentials.json found\n"+
				"To configure google-driver:\n"+
				"1. Follow the GCP setup guide in outlines/google_cloud_platform_setup.md\n"+
				"2. Download your credentials.json from Google Cloud Console\n"+
				"3. Save it to: %s\n",
			filepath.Join(homeDir, authCredentialPath),
		)
	}

	return google.ConfigFromJSON(data, drive.DriveScope)

}

func PrintConfig(c *oauth2.Config) {
	fmt.Printf("Client ID: %v\n", c.ClientID)
	fmt.Printf("Client Secret: %v\n", c.ClientSecret)
	fmt.Printf("Client Endpoint: %v\n", c.Endpoint)
	fmt.Printf("Redirect URL: %v\n", c.RedirectURL)
	fmt.Printf("Scopes: %v\n", c.Scopes)
}
