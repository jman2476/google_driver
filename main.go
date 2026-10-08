package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

type apiConfig struct {
	client  *http.Client
	token   *oauth2.Token
	service *drive.Service
}

type arguments struct {
	tool   string
	source string
	target string
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
		fmt.Printf("Succesfully logged in")
		os.Exit(0)
	case "logout":
		_, err := config.commandLogOut(args.source, args.target)
		if err != nil {
			fmt.Printf("Log out error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Succesfully logged out")
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

	// about, err := config.service.About.Get().Fields("user").Do()
	// if err != nil {
	// 	log.Fatalf("Drive API check failed: %v", err)
	// }

	// fmt.Printf(
	// 	"Authenticated successfully as %s (%s)\n",
	// 	about.User.DisplayName,
	// 	about.User.EmailAddress,
	// )
}
