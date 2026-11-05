package argparser

import (
	"fmt"
	"strings"
)

type FlagArg struct {
	flag  string
	value string
}

type Arguments struct {
	Command    string
	SourcePath string
	TargetPath string
	Flags      []FlagArg
}

func ParseFlags(args []string) (parsed Arguments) {
	for i, a := range args {
		fmt.Println(i, a)
	}

	parsed.Command = args[1]
	parsed.SourcePath = args[2]
	parsed.TargetPath = args[3]

	if len(args) > 4 {
		for i := 4; i < len(args); i++ {
			if strings.Contains(args[i], "-") {
				newflag := FlagArg{flag: args[i]}

				if i+1 < len(args) && !strings.Contains(args[i+1], "-") {
					newflag.value = args[i+1]
					i++
				}
				parsed.Flags = append(parsed.Flags, newflag)
			}
		}
	}

	PrintParsed(parsed)

	return
}

func PrintParsed(parsed Arguments) {
	fmt.Printf("Command: %s\n", parsed.Command)
	fmt.Printf("Source Path: %s\n", parsed.SourcePath)
	fmt.Printf("Target Path: %s\n", parsed.TargetPath)
	fmt.Println("Flags:")

	for _, f := range parsed.Flags {
		fmt.Println(f)
	}
}
