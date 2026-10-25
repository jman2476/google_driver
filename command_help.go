package main

import (
	"fmt"
	"strings"

	"google.golang.org/api/drive/v3"
)

func (cfg *apiConfig) commandHelp(command, _ string) (*drive.File, error) {
	helpMessage := "call 'help <command name>' to get information on a specific command\n"
	var msgWriter strings.Builder
	msgWriter.Write([]byte(helpMessage))

	if command != "" {
		msgWriter.Write([]byte(
			formatHelp(command, cfg.commandRegistry()[command]),
		))
	} else {
		for call, cmd := range cfg.commandRegistry() {
			msgWriter.Write([]byte(
				formatHelp(call, cmd),
			))
		}
	}

	fmt.Println(msgWriter.String())
	return nil, nil
}

func formatHelp(call string, command cliCommand) string {
	var strBuilder strings.Builder
	strBuilder.Write(
		fmt.Appendf(
			[]byte{},
			"%s --> %s\n     %s\n",
			call, command.name,
			command.description,
		),
	)

	return strBuilder.String()
}
