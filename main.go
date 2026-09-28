package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

type apiConfig struct {
	client  *http.Client
	token   *oauth2.Token
	service *drive.Service
}

func main() {
	var config apiConfig
	err := config.setClient()
	if err != nil {
		log.Fatalf("Error setting client: %v\n", err)
	}

	err = config.setService()
	if err != nil {
		log.Fatalf("Error setting service: %v\n", err)
	}

	log.Println("Google Drive service initialized")

	about, err := config.service.About.Get().Fields("user").Do()
	if err != nil {
		log.Fatalf("Drive API check failed: %v", err)
	}

	fmt.Printf(
		"Authenticated successfully as %s (%s)\n",
		about.User.DisplayName,
		about.User.EmailAddress,
	)
}
