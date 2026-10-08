package main

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

func (cfg *apiConfig) commandAccount(_, _ string) (*drive.File, error) {
	about, err := cfg.service.About.Get().Fields("user").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to get current account: %w", err)
	}

	fmt.Printf(
		"Logged in as %s (%s)\n",
		about.User.DisplayName,
		about.User.EmailAddress,
	)

	return nil, nil
}
