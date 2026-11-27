package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"google.golang.org/api/drive/v3"
)

var (
	ErrNeedRecursive       = errors.New("to download a folder, use the '-r' flag")
	ErrFileNotFound        = errors.New("no file found at this address")
	ErrDuplicateFilesFound = errors.New("multiple files found w/ same name")
)

func (cfg *apiConfig) downloadFile(path string, targetPath string, flags map[string]string) (*drive.File, error) {
	dirId, filename, isFolder, err := cfg.ValidateDrivePath(path)
	if err != nil {
		return nil, err
	}

	if isFolder {
		if _, ok := flags["r"]; ok {
			return cfg.downloadFolder(path, targetPath, flags)
		}
		return nil, ErrNeedRecursive
	}
	query := fmt.Sprintf(
		"name = '%s' and '%s' in parents and trashed = false",
		filename,
		dirId,
	)

	files, err := cfg.service.Files.List().
		Q(query).
		Fields("files(id, name)").
		Do()

	if err != nil {
		return nil, err
	}

	if len(files.Files) == 0 {
		return nil, ErrFileNotFound
	}
	if len(files.Files) > 1 {
		cont, err := cfg.multiFilePrompt(filename)
		if err != nil || !cont {
			return nil, ErrDuplicateFilesFound
		}
	}

	fileId := files.Files[0].Id
	res, err := cfg.service.Files.Get(fileId).Download()
	if err != nil {
		return nil, fmt.Errorf(
			"Error downloading file: Status %d; %s\n",
			res.StatusCode, err,
		)
	}
	defer res.Body.Close()

	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("cannot read current directory: %w", err)
	}
	defer os.Chdir(pwd)

	err = os.Chdir(targetPath)
	if err != nil {
		fmt.Printf(
			"Target directory %s not found, writing to %s instead",
			targetPath, pwd,
		)
	}

	currDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf(
			"cannot read current directory after deciding download location: %w", err,
		)
	}

	destPath := filepath.Join(currDir, filename)
	outFile, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create local file at %s: %w",
			destPath, err,
		)
	}
	defer outFile.Close()

	written, err := io.Copy(outFile, res.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to write to local file at %s: %w",
			destPath, err,
		)
	}

	fmt.Printf(
		"Downloaded '%s' (%d bytes)  to '%s'\n",
		filename,
		written,
		destPath,
	)

	return nil, nil
}

func (cfg *apiConfig) downloadFolder(path string, targetPath string, flags map[string]string) (*drive.File, error) {

	return nil, nil
}

func (cfg *apiConfig) multiFilePrompt(filename string) (bool, error) {
	prompt := fmt.Sprintf("Multiple instances of file '%s' found. Download the first one? (y/n)    ", filename)
	input := bufio.NewReader(cfg.in)
	scanner := bufio.NewScanner(input)
	fmt.Print(prompt)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf(
				"error reading input: %w", err,
			)
		}

		return false, ErrScannerEOF
	}
	result := scanner.Text()

	if result != "y" && result != "Y" {
		fmt.Println("Canceling operation.")
		return false, nil
	}

	return true, nil
}
