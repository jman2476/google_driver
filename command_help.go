package main

import (
	"fmt"
	"strings"

	"google.golang.org/api/drive/v3"
)

func (cfg *apiConfig) commandHelp(command, _ string) (*drive.File, error) {
	var helpMessage string

	if command != "" {
		helpMessage += formatHelp(command, cfg.commandRegistry()[command])
	} else {
		for call, cmd := range cfg.commandRegistry() {
			helpMessage += formatHelp(call, cmd)
		}
	}

	fmt.Println(helpMessage)
	return nil, nil
}

func formatHelp(call string, command cliCommand) string {
	var strBuilder strings.Builder
	strBuilder.Write(
		fmt.Appendf(
			[]byte{},
			"%s --> %s\n     %s",
			call, command.name,
			command.description,
		),
	)

	return strBuilder.String()
}
