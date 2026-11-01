package main

import "google.golang.org/api/drive/v3"

type cliCommand struct {
	name        string
	use         string
	description string
	callback    func(string, string) (*drive.File, error)
}

func (cfg *apiConfig) commandRegistry() map[string]cliCommand {
	return map[string]cliCommand{
		"ls": {
			name:     "List directory",
			use:      "List contents of the given directory",
			callback: cfg.commandList,
		},
		"rm": {
			name:     "Remove",
			use:      "Remove/delete a file or directory",
			callback: cfg.commandRemove,
		},
		"mv": {
			name:        "Move",
			use:         "Upload file to your Google Drive",
			description: "mv <file source> <drive target>\n    *drive target can be a folder on your drive and/or a new file name. If target folder is omitted, file will be uploaded to root directory on drive",
			callback:    cfg.uploadFile,
		},
		"mkdir": {
			name:     "Make directory",
			use:      "make new directory on your Google Drive",
			callback: cfg.createFolder,
		},
		"login": {
			name:     "Log In",
			use:      "Log into your Google account",
			callback: cfg.commandLogIn,
		},
		"logout": {
			name:     "Log Out",
			use:      "Log out of your Google account",
			callback: cfg.commandLogOut,
		},
		"account": {
			name:     "Account details",
			use:      "Get current account",
			callback: cfg.commandAccount,
		},
		"help": {
			name:     "Help menu",
			use:      "Display help information",
			callback: cfg.commandHelp,
		},
	}
}
