package main

import (
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

		if len(fileData.Data) != int(fileData.Info.Size()) {
			t.Errorf(
				"Fail: size mismatch\nExpected: %v [.Size()]\nActual: %v [len(Data)]",
				int(fileData.Info.Size()), len(fileData.Data),
			)
		}
	}
}
