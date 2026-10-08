package main

import "google.golang.org/api/drive/v3"

type cliCommand struct {
	name        string
	description string
	callback    func(string, string) (*drive.File, error)
}

func (cfg *apiConfig) commandRegister(command string) map[string]cliCommand {
	return map[string]cliCommand{
		"rm": {},
		"mv": {},
		"mkdir": {
			name:        "make directory on drive",
			description: "",
			callback:    cfg.createFolder,
		},
		"login": {
			name:        "Log In",
			description: "Log into your Google account",
			callback:    cfg.commandLogIn,
		},
		"logout": {
			name:        "Log Out",
			description: "Log out of your Google account",
			callback:    cfg.commandLogOut,
		},
		"account": {
			name:        "Account details",
			description: "Get current account",
			callback:    cfg.commandAccount,
		},
	}
}
