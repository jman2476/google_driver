package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

type apiConfig struct {
	client  *http.Client
	token   *oauth2.Token
	service *drive.Service
	in      io.Reader
}

func main() {
	var config apiConfig

	args, err := parseArgs(os.Args)
	if err != nil {
		fmt.Println("Argument error:", err)
		os.Exit(1)
	}

	switch args.tool {
	case "login":
		_, err := config.commandLogIn(args.source, args.target)
		if err != nil {
			fmt.Printf("Log in error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully logged in\n")
		os.Exit(0)
	case "logout":
		_, err := config.commandLogOut(args.source, args.target)
		if err != nil {
			fmt.Printf("Log out error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully logged out\n")
		os.Exit(0)
	default:
		err := config.setClient()
		if err != nil {
			fmt.Printf("Error setting client: %v\n", err)
			os.Exit(1)
		}

		err = config.setService()
		if err != nil {
			fmt.Printf("Error setting service: %v\n", err)
			os.Exit(1)
		}

		config.in = os.Stdin

		fmt.Println("Google Drive service initialized")
		tool := config.commandRegistry()[args.tool]
		file, err := tool.callback(args.source, args.target)
		if err != nil {
			fmt.Printf("%s error: %v\n", tool.name, err)
			os.Exit(1)
		}
		if file != nil {
			fmt.Printf(
				"Performed '%s' on '%s'\n",
				tool.name, file.Name,
			)
		}
		os.Exit(0)
	}

}
