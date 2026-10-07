package main

type cliCommand struct {
	name        string
	description string
	callback    func(*apiConfig) error
}

func commandRegister(command string) map[string]cliCommand {
	return map[string]cliCommand{
		"rm":     {},
		"mv":     {},
		"mkdir":  {},
		"login":  {},
		"logout": {},
	}
}
