package main

import (
	"fmt"
	"log"

	"github.com/jman2476/google-driver/internal/auth"
)

func main() {
	log.Println("Starting Google Driver")
	config, err := auth.LoadCredentials()
	if err != nil {
		fmt.Printf("Error loading credentials: %v", err)
	}

	auth.PrintConfig(config)

	token, err := auth.ReadTokenData()
	if err != nil {
		fmt.Printf("Token error: %v\n", err)

		authCode, err := auth.GetAuthorization(config)
		if err != nil {
			fmt.Printf("Error getting authorization: %v", err)
		} else {
			token, err := auth.GetToken(config, authCode)
			if err != nil {
				fmt.Printf("Error getting token: %v", err)
			} else {
				fmt.Printf("Token pointer: %v", token)
				err = auth.WriteTokenData(token)
				if err != nil {
					fmt.Printf("Error writing token cache: %v", err)
				} else {
					fmt.Printf("Token cache written")
				}
			}

		}
	} else {
		fmt.Printf("Token: %v\n", token)
	}

}
