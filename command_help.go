package main

import (
	"fmt"
	"strings"

	"google.golang.org/api/drive/v3"
)

func (cfg *apiConfig) commandHelp(args arguments) (*drive.File, error) {
	helpMessage := "call 'help <command name>' to get information on a specific command\n"
	var msgWriter strings.Builder
	msgWriter.Write([]byte(helpMessage))

	if len(args.positional) == 0 {
		for call, cmd := range cfg.commandRegistry() {
			msgWriter.Write([]byte(
				formatHelp(call, cmd),
			))
		}
	} else {
		command := args.positional[0]
		msgWriter.Write([]byte(
			formatHelpFull(command, cfg.commandRegistry()[command]),
		))
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
			command.use,
		),
	)

	return strBuilder.String()
}

func formatHelpFull(call string, command cliCommand) string {
	var msgBuilder strings.Builder
	msgBuilder.Write([]byte(formatHelp(call, command)))
	msgBuilder.Write(
		fmt.Appendf(
			[]byte{},
			"\n%s\n",
			command.description,
		),
	)

	return msgBuilder.String()
}
