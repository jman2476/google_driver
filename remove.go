package main

import (
	"bufio"
	"errors"
	"fmt"
	"slices"

	"google.golang.org/api/drive/v3"
)

var (
	ErrFileNotFound      = errors.New("File delete: file not found")
	ErrTooManyFilesFound = errors.New(
		"File delete: more than one file found with same name",
	)
	ErrRemoveCanceled = errors.New("File delete: user canceled")
	ErrScannerEOF     = errors.New("File delete: scanner hit EOF")
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
		Fields("files(id, name, mimeType)").
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

	var toTrash bool
	input := bufio.NewReader(cfg.in)
	scanner := bufio.NewScanner(input)
	fileType := "file"
	if res.Files[0].MimeType == driveFolderMIME {
		fileType = "folder"
	}
	// message := fmt.Sprintf("Are you sure you want to delete this %s? (y/n)  ", fileType)
	userPrompts := []string{
		fmt.Sprintf("Are you sure you want to delete this %s?", fileType),
		"If you delete this folder, all contents will be deleted as well. Continue?",
		"Move to trash instead?",
	}

	for i := 0; i < 3; {
		if i == 1 && fileType != "folder" {
			i++
			continue
		}
		fmt.Printf("%s (y/n)   ", userPrompts[i])
		if !scanner.Scan() {

			if err = scanner.Err(); err != nil {
				return nil, fmt.Errorf(
					"error reading input: %w", err,
				)
			}

			return nil, ErrScannerEOF
		}

		result := scanner.Text()
		if len(result) > 1 ||
			len(result) < 1 ||
			!slices.Contains([]string{"y", "n", "Y", "N"}, result) {
			fmt.Println("please reply (y/n)")
			continue
		} else if slices.Contains([]string{"n", "N"}, result) &&
			i != 2 {
			return nil, ErrRemoveCanceled
		} else {
			if slices.Contains([]string{"n", "N"}, result) &&
				i == 2 {
				toTrash = true
			}
			i++
		}

	}

	// for valid := true; valid; {
	// 	fmt.Print(message)
	// 	scanner.Scan()

	// 	if err = scanner.Err(); err != nil {
	// 		return nil, fmt.Errorf(
	// 			"error reading input: %w", err,
	// 		)
	// 	}

	// 	result := scanner.Text()
	// 	if len(result) > 1 ||
	// 		len(result) < 1 ||
	// 		!slices.Contains([]string{"y", "n", "Y", "N"}, result) {
	// 		fmt.Println("please reply (y/n)")
	// 		continue
	// 	} else if slices.Contains([]string{"n", "N"}, result) {
	// 		return nil, ErrRemoveCanceled
	// 	} else {
	// 		valid = false
	// 	}
	// }
	// if fileType == "folder" {
	// 	message = "If you delete this folder, all contents will be deleted as well. Continue? (y/n)   "

	// 	for valid := true; valid; {
	// 		fmt.Print(message)
	// 		scanner.Scan()

	// 		if err = scanner.Err(); err != nil {
	// 			return nil, fmt.Errorf(
	// 				"error reading input: %w", err,
	// 			)
	// 		}

	// 		result := scanner.Text()
	// 		if len(result) > 1 ||
	// 			len(result) < 1 ||
	// 			!slices.Contains([]string{"y", "n", "Y", "N"}, result) {
	// 			fmt.Println("please reply (y/n)")
	// 			continue
	// 		} else if slices.Contains([]string{"n", "N"}, result) {
	// 			return nil, ErrRemoveCanceled
	// 		} else {
	// 			valid = false
	// 		}
	// 	}
	//}

	var file *drive.File
	switch toTrash {
	case true:
		file, err = cfg.service.Files.
			Update(res.Files[0].Id, &drive.File{Trashed: true}).Do()
		if err != nil {
			return nil, fmt.Errorf(
				"unable to delete file: %w", err,
			)
		}
		return file, nil
	case false:
		err = cfg.service.Files.Delete(res.Files[0].Id).Do()
		if err != nil {
			return nil, fmt.Errorf(
				"unable to delete file: %w", err,
			)
		}
		file = &drive.File{Name: fileName}
	}

	return file, nil
}
