package main

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

func (cfg *apiConfig) commandRemove(_, _ string) (*drive.File, error) {
	return nil, fmt.Errorf("Not implemented")
}
