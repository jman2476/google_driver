package main

import (
	"errors"
)

var (
	ErrMissingArgs = errors.New("not enough arguments")
	ErrTooManyArgs = errors.New("too many arguments")
)

func parseArgs(args []string) (
	parsed arguments, err error) {
	switch len(args) {
	case 1:
		err = ErrMissingArgs
	case 2:
		parsed.tool = args[1]
	case 3:
		parsed.tool = args[1]
		parsed.source = args[2]
	case 4:
		parsed.tool = args[1]
		parsed.source = args[2]
		parsed.target = args[3]
	default:
		err = ErrTooManyArgs
	}

	return
}
