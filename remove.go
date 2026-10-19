package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"

	"google.golang.org/api/drive/v3"
)

var (
	ErrFileNotFound      = errors.New("File delete: file not found")
	ErrTooManyFilesFound = errors.New(
		"File delete: more than one file found with same name",
	)
	ErrRemoveCanceled = errors.New("File delete: user canceled")
)

func (cfg *apiConfig) removeFile(fileName, targetDir string) (*drive.File, error) {
	location, err := cfg.findFolder(targetDir)
	if err != nil {
		return nil, fmt.Errorf("unable to locate target directory: %w", err)
	}

	query := fmt.Sprintf(
		"name = '%s' and '%s' in parents and trashed = false",
		fileName, location.Id,
	)

	res, err := cfg.service.Files.List().
		Q(query).
		Fields("files(id, name)").
		Do()

	if err != nil {
		return nil, fmt.Errorf(
			"error finding file '%s': %w",
			fileName, err,
		)
	}

	if len(res.Files) == 0 {
		return nil, fmt.Errorf(
			"%w: %s", ErrFileNotFound, fileName,
		)
	}

	if len(res.Files) > 1 {
		return nil, fmt.Errorf(
			"%w: %s", ErrTooManyFilesFound, fileName,
		)
	}

	fmt.Printf("Found file: %v\n", res.Files[0].Name)

	input := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(input)

	for valid := true; valid; {
		fmt.Print("Are you sure you want to delete this file? (y/n)  ")
		scanner.Scan()

		if err = scanner.Err(); err != nil {
			return nil, fmt.Errorf(
				"error reading input: %w", err,
			)
		}

		result := scanner.Text()
		if len(result) > 1 ||
			len(result) < 1 ||
			!slices.Contains([]string{"y", "n", "Y", "N"}, result) {
			fmt.Println("please reply (y/n)")
			continue
		} else if slices.Contains([]string{"n", "N"}, result) {
			return nil, ErrRemoveCanceled
		} else {
			valid = false
		}
	}

	err = cfg.service.Files.Delete(res.Files[0].Id).Do()
	if err != nil {
		return nil, fmt.Errorf(
			"unable to delete file: %w", err,
		)
	}
	return &drive.File{Name: fileName}, nil
}
