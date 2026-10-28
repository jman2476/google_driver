package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindFile(t *testing.T) {
	currentDir, err := filepath.Abs(".")
	if err != nil {
		t.Errorf("Error: cannot get working directory: %v", err)
		return
	}

	cases := []struct {
		inputPath        string
		expectedPath     string
		expectedFileName string
		exists           bool
	}{
		{
			inputPath: "./test_resources/goose_sqlc-instructions.txt",
			expectedPath: filepath.Join(
				currentDir,
				"./test_resources/goose_sqlc-instructions.txt"),
			expectedFileName: "goose_sqlc-instructions.txt",
			exists:           true,
		}, {
			inputPath:        "/test_resources/goose_sqlc-instructions.txt",
			expectedPath:     "/test_resources/goose_sqlc-instructions.txt",
			expectedFileName: "goose_sqlc_instructions.txt",
			exists:           false,
		}, {
			inputPath: "./test_resources/goose_sqlc-instructions.md",
			expectedPath: filepath.Join(
				currentDir,
				"./test_resources/goose_sqlc-instructions.md",
			),
			expectedFileName: "goose_sqlc-instructions.md",
			exists:           false,
		}, {
			inputPath: "./test_resources/bubbletea_gui.gif",
			expectedPath: filepath.Join(
				currentDir,
				"./test_resources/bubbletea_gui.gif",
			),
			expectedFileName: "bubbletea_gui.gif",
			exists:           true,
		}, {
			inputPath: "./test_resources/T-Rex_training.cpp",
			expectedPath: filepath.Join(
				currentDir,
				"./test_resources/T-Rex_training.cpp",
			),
			expectedFileName: "T-Rex_training.cpp",
			exists:           false,
		}, {
			inputPath: "./test_resources/OrionNebulaCenterMonitor.png",
			expectedPath: filepath.Join(
				currentDir,
				"./test_resources/OrionNebulaCenterMonitor.png",
			),
			expectedFileName: "OrionNebulaCenterMonitor.png",
			exists:           true,
		}, {
			inputPath: "./test_resurces/goose_sqlc_instructions.txt",
			expectedPath: filepath.Join(
				currentDir,
				"./test_resurces/goose_sqlc_instructions.txt",
			),
			expectedFileName: "goose_sqlc_instructions.txt",
			exists:           false,
		}, {
			inputPath: "./test_resources/chess.mov",
			expectedPath: filepath.Join(
				currentDir,
				"./test_resources/chess.mov",
			),
			expectedFileName: "chess.mov",
			exists:           true,
		},
	}

	for _, c := range cases {
		resultPath, resultInfo, err := FindFile(c.inputPath)

		if err != nil && c.exists {
			t.Errorf(
				"Fail: error when finding existing file at %s: %v",
				c.expectedPath, err,
			)
		} else if c.exists {
			if resultPath != c.expectedPath {
				t.Errorf(
					"Fail: Paths don't match: \nExpected: %s\nActual: %s",
					c.expectedPath, resultPath,
				)
			}
			if resultInfo.Name() != c.expectedFileName {
				t.Errorf(
					"Fail: Names don't match: \nExpected: %s\nActual: %s",
					c.expectedFileName, resultInfo.Name(),
				)
			}
		}
	}
}

func TestValidateFile(t *testing.T) {
	currentDir, err := filepath.Abs(".")
	if err != nil {
		t.Errorf("Error: cannot get working directory: %v", err)
		return
	}

	cases := []struct {
		path          string
		expectAbsPath string
		expectName    string
		expectMime    string
	}{
		{
			path: "./test_resources/bubbletea_gui.gif",
			expectAbsPath: filepath.Join(
				currentDir, "./test_resources/bubbletea_gui.gif",
			),
			expectName: "bubbletea_gui.gif",
			expectMime: "image/gif",
		}, {
			path: "./test_resources/chess.mov",
			expectAbsPath: filepath.Join(
				currentDir, "./test_resources/chess.mov",
			),
			expectName: "chess.mov",
			expectMime: "video/quicktime",
		}, {
			path: "./test_resources/OrionNebulaCenterMonitor.png",
			expectAbsPath: filepath.Join(
				currentDir,
				"./test_resources/OrionNebulaCenterMonitor.png",
			),
			expectName: "OrionNebulaCenterMonitor.png",
			expectMime: "image/png",
		}, {
			path: "./test_resources/goose_sqlc-instructions.txt",
			expectAbsPath: filepath.Join(
				currentDir,
				"./test_resources/goose_sqlc-instructions.txt",
			),
			expectName: "goose_sqlc-instructions.txt",
			expectMime: "text/plain",
		}, {
			path: "./test_resources/LICENSE",
			expectAbsPath: filepath.Join(
				currentDir,
				"./test_resources/LICENSE",
			),
			expectName: "LICENSE",
			expectMime: "text/plain",
		}, {
			path: "./test_resources/NO_EXTENSION",
			expectAbsPath: filepath.Join(
				currentDir,
				"./test_resources/NO_EXTENSION",
			),
			expectName: "NO_EXTENSION",
			expectMime: "text/plain",
		},
	}

	for _, c := range cases {
		fileData, err := ValidateFile(c.path)
		if err != nil {
			t.Errorf("Fail: unable to validate file: %v", err)
		}

		if fileData.Info.Name() != c.expectName {
			t.Errorf(
				"Fail: name mismatch\nExpected: %v\nActual: %v",
				c.expectName, fileData.Info.Name(),
			)
		}

		if strings.Split(fileData.Mime, ";")[0] != c.expectMime {
			t.Errorf(
				"Fail: mime type mismatch\nExpected: %v\nActual: %v",
				c.expectMime, fileData.Mime,
			)
		}

		if fileData.AbsPath != c.expectAbsPath {
			t.Errorf(
				"Fail: absolute path mismatch\nExpected: %v\nActual: %v",
				c.expectAbsPath, fileData.AbsPath,
			)
		}

		t.Logf("Size of file: %d bytes\n", int(fileData.Info.Size()))
		// t.Logf("Len of data slice: %d \n", len(fileData.Data))

		// if len(fileData.Data) != int(fileData.Info.Size()) {
		// 	t.Errorf(
		// 		"Fail: size mismatch\nExpected: %v [.Size()]\nActual: %v [len(Data)]",
		// 		int(fileData.Info.Size()), len(fileData.Data),
		// 	)
		// }
	}
}

func TestValidateUnreadableFile(t *testing.T) {
	// If root, skip test
	if os.Geteuid() == 0 {
		t.Skip("Skipping permission test, running as root")
	}

	tempDir := t.TempDir()
	unreadablePath := filepath.Join(
		tempDir, "unreadable.txt")

	err := os.WriteFile(
		unreadablePath,
		[]byte("Testing text for texty testing"),
		0600,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create temp test file: %v", err)
	}

	err = os.Chmod(unreadablePath, 0000)
	if err != nil {
		t.Fatalf(
			"Failed to set file permissions to unreadable: %v",
			err,
		)
	}

	t.Cleanup(func() {
		_ = os.Chmod(unreadablePath, 0600)
	})

	fileData, err := ValidateFile(unreadablePath)
	if err == nil {
		t.Errorf(
			"Fail: Validate file returned nil error on the unreadable file: %v",
			fileData,
		)
	}

}

func TestValidateUploadPathRoot(t *testing.T) {
	var config apiConfig
	err := config.setClient()
	if err != nil {
		log.Fatalf("Error setting client: %v\n", err)
	}

	err = config.setService()
	if err != nil {
		log.Fatalf("Error setting service: %v\n", err)
	}

	cases := []string{
		"",
		"root",
		"/root",
		"root/",
		"/root/",
	}

	for _, c := range cases {
		dirID, name, isFolder, err := config.ValidateUploadPath(c)

		if dirID != "root" {
			t.Errorf(
				"Fail: case: %s\nExpected: %s\n Actual: %s",
				c, "root", dirID,
			)
		}
		if name != "" {
			t.Errorf(
				"Fail: case: %s\nExpected: [%s]\n Actual: [%s]",
				c, "", name,
			)
		}
		if !isFolder {
			t.Errorf(
				"Fail: case: %s\nExpected: true \n Actual: %t",
				c, isFolder,
			)
		}
		if err != nil {
			t.Errorf(
				"Fail: case: %s\nExpected: nil\n Actual: %v",
				c, err,
			)
		}
	}
}

func TestValidateUploadPathIsFolder(t *testing.T) {
	var config apiConfig
	err := config.setClient()
	if err != nil {
		log.Fatalf("Error setting client: %v\n", err)
	}

	err = config.setService()
	if err != nil {
		log.Fatalf("Error setting service: %v\n", err)
	}

	type expected struct {
		dirId    string
		name     string
		isFolder bool
		err      error
	}

	cases := []struct {
		path string
		expected
	}{
		{
			path: "/test_folder/",
			expected: expected{
				dirId:    "1qXisdgvsZx1V-YTyIWp0Ha28HOGzw5WF",
				name:     "",
				isFolder: true,
				err:      nil,
			},
		}, {
			path: "test_folder",
			expected: expected{
				dirId:    "1qXisdgvsZx1V-YTyIWp0Ha28HOGzw5WF",
				name:     "",
				isFolder: true,
				err:      nil,
			},
		}, {
			path: "test_folder/child_test_folder",
			expected: expected{
				dirId:    "1S8x3ni_BVO5H7usudX3RJ1b79jb5P4Fi",
				name:     "",
				isFolder: true,
				err:      nil,
			},
		}, {
			path: "test_folder/boogie_wonderland/",
			expected: expected{
				dirId:    "",
				name:     "",
				isFolder: true,
				err:      ErrFolderDNE,
			},
		}, {
			path: "test_folder/child_test_folder/chess.mov/",
			expected: expected{
				dirId:    "",
				name:     "",
				isFolder: true,
				err:      ErrFolderDNE,
			},
		}, {
			path: "test_folder/child_test_folder/chess.mov/",
			expected: expected{
				dirId:    "",
				name:     "",
				isFolder: true,
				err:      ErrFolderDNE,
			},
		}, {
			path: "test_folder/child_test_folder/dot.test.folder",
			expected: expected{
				dirId:    "13Ao5zWzDw38A4WKhF1zOPrC3tRnbs4UV",
				name:     "",
				isFolder: true,
				err:      nil,
			},
		}, {
			path: "test_folder/child_test_folder/dot.folder_test",
			expected: expected{
				dirId:    "1gUS-RdhGC1Vojna3E3b1SZFq6UpFs6om",
				name:     "",
				isFolder: true,
				err:      nil,
			},
		},
	}

	for _, c := range cases {
		dirID, name, isFolder, err := config.ValidateUploadPath(c.path)
		expectedErr := fmt.Errorf("%v: %s", ErrFolderDNE, c.path)
		if err == nil {
			expectedErr = err
		}
		if dirID != c.dirId {
			t.Errorf(
				"Fail: case: %s\nExpected: %s\n Actual: %s",
				c.path, c.dirId, dirID,
			)
		}
		if name != c.name {
			t.Errorf(
				"Fail: case: %s\nExpected: [%s]\n Actual: [%s]",
				c.path, c.name, name,
			)
		}
		if isFolder != c.isFolder {
			t.Errorf(
				"Fail: case: %s\nExpected: %t \n Actual: %t",
				c.path, c.isFolder, isFolder,
			)
		}
		if err != c.err {
			if c.err == nil && err == nil {
				continue
			}
			if !errors.Is(err, c.err) {
				t.Errorf(
					"Fail: case: %s\nExpected: %v\n Actual: %v",
					c.path,
					expectedErr,
					err,
				)
			}
		}
	}
}

func TestValidateUploadPathIsFile(t *testing.T) {
	var config apiConfig
	err := config.setClient()
	if err != nil {
		log.Fatalf("Error setting client: %v\n", err)
	}

	err = config.setService()
	if err != nil {
		log.Fatalf("Error setting service: %v\n", err)
	}

	type expected struct {
		dirId    string
		name     string
		isFolder bool
		err      error
	}

	cases := []struct {
		path string
		expected
	}{
		{
			path: "test_folder/bingo.txt",
			expected: expected{
				dirId:    "1qXisdgvsZx1V-YTyIWp0Ha28HOGzw5WF",
				name:     "bingo.txt",
				isFolder: false,
				err:      nil,
			},
		}, {
			path: "how-to-brush-teeth.txt",
			expected: expected{
				dirId:    "root",
				name:     "how-to-brush-teeth.txt",
				isFolder: false,
				err:      nil,
			},
		}, {
			path: "how-to-brush-teeth.txt/",
			expected: expected{
				dirId:    "",
				name:     "",
				isFolder: true,
				err:      ErrFolderDNE,
			},
		},
	}

	for _, c := range cases {
		dirID, name, isFolder, err := config.ValidateUploadPath(c.path)
		expectedErr := fmt.Errorf("%v: %s", ErrFolderDNE, c.path)
		if err == nil {
			expectedErr = err
		}
		if dirID != c.dirId {
			t.Errorf(
				"Fail: case: %s\nExpected: %s\n Actual: %s",
				c.path, c.dirId, dirID,
			)
		}
		if name != c.name {
			t.Errorf(
				"Fail: case: %s\nExpected: [%s]\n Actual: [%s]",
				c.path, c.name, name,
			)
		}
		if isFolder != c.isFolder {
			t.Errorf(
				"Fail: case: %s\nExpected: %t \n Actual: %t",
				c.path, c.isFolder, isFolder,
			)
		}
		if err != c.err {
			if c.err == nil && err == nil {
				continue
			}
			if !errors.Is(err, c.err) {
				t.Errorf(
					"Fail: case: %s\nExpected: %v\n Actual: %v",
					c.path,
					expectedErr,
					err,
				)
			}
		}
	}

}
