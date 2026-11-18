package main

import (
	"errors"
	"fmt"
	"strings"
)

type arguments struct {
	tool       string
	positional []string
	flags      map[string]string
}

var (
	ErrMissingCmd  = errors.New("missing command")
	ErrMissingArgs = errors.New("not enough arguments")
	ErrTooManyArgs = errors.New("too many arguments")
)

func (a arguments) HasFlag(names ...string) bool {
	for _, name := range names {
		if _, ok := a.flags[name]; ok {
			return true
		}
	}

	return false
}

func (a arguments) GetFlag(names ...string) string {
	for _, name := range names {
		if val, ok := a.flags[name]; ok {
			return val
		}
	}
	return ""
}

func parseArgs(rawArgs []string) (arguments, error) {
	printParse(rawArgs)
	if len(rawArgs) < 2 {
		return arguments{}, ErrMissingCmd
	}

	parsed := arguments{
		tool:  rawArgs[1],
		flags: make(map[string]string),
	}

	args := rawArgs[2:]
	for i := 0; i < len(args); i++ {
		token := args[i]

		if strings.HasPrefix(token, "-") {
			flagName := strings.TrimLeft(token, "-")

			if strings.Contains(flagName, "=") {
				parts := strings.SplitN(flagName, "=", 2)
				parsed.flags[parts[0]] = parts[1]
				continue
			}

			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				parsed.flags[flagName] = args[i+1]
			} else {
				parsed.flags[flagName] = "true"
			}
			continue
		}

		parsed.positional = append(parsed.positional, token)
	}

	return parsed, nil
}

func printParse(args []string) {
	for i, a := range args {
		fmt.Printf("%d: %s\n", i, a)
	}
}
