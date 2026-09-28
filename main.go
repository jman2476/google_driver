package main

import (
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
}
