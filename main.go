package main

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type apiConfig struct {
	client *http.Client
	token  *oauth2.Token
}

func main() {
	var config apiConfig
	err := config.setClient()
	if err != nil {
		fmt.Printf("Error setting client: %v\n", err)
	} else {
		fmt.Printf("Client set: %v\n", config)
	}

}
