package main

import (
	"errors"
	"fmt"
)

type arguments struct {
	tool       string
	source     string
	target     string
	positional []string
	flags      map[string]string
}

var (
	ErrMissingCmd  = errors.New("missing command")
	ErrMissingArgs = errors.New("not enough arguments")
	ErrTooManyArgs = errors.New("too many arguments")
)

func parseArgs(rawArgs []string) (
	parsed arguments, err error) {
	printParse(rawArgs)
	if len(rawArgs) < 2 {
		return parsed, ErrMissingCmd
	}

	switch len(rawArgs) {
	case 2:
		parsed.tool = rawArgs[1]
	case 3:
		parsed.tool = rawArgs[1]
		parsed.source = rawArgs[2]
	case 4:
		parsed.tool = rawArgs[1]
		parsed.source = rawArgs[2]
		parsed.target = rawArgs[3]
	default:
		err = ErrTooManyArgs
	}

	return
}

func printParse(args []string) {
	for i, a := range args {
		fmt.Printf("%d: %s\n", i, a)
	}
}
