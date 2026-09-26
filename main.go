package main

import (
	"fmt"
	"log"

	"github.com/jman2476/google-driver/internal/auth"
)

func main() {
	config, err := auth.LoadCredentials()
	if err != nil {
		fmt.Printf("Error loading credentials: %v", err)
	}

	auth.PrintConfig(config)

	token, err := auth.ReadTokenData()
	if err != nil {
		fmt.Printf("Token error: %v\n", err)

		auth.GetAuthorization(config)
	} else {
		fmt.Printf("Token: %v\n", token)
	}

	log.Println("Starting Google Driver")
}
