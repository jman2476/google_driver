package main

import (
	"log"
	"testing"

	"google.golang.org/api/drive/v3"
)

func TestWriteUploadMetadata(t *testing.T) {
	cases := []struct {
		path         string
		targetFolder string
		expected     *drive.File
	}{
		{
			path:         "./test_resources/bubbletea_gui.gif",
			targetFolder: "test_folder",
			expected: &drive.File{
				Name:     "bubbletea_gui.gif",
				MimeType: "image/gif",
				Parents:  []string{"test_folder"},
			},
		}, {
			path:         "./test_resources/chess.mov",
			targetFolder: "",
			expected: &drive.File{
				Name:     "chess.mov",
				MimeType: "video/quicktime",
			},
		}, {
			path:         "./test_resources/goose_sqlc-instructions.txt",
			targetFolder: "child_test_folder",
			expected: &drive.File{
				Name:     "goose_sqlc-instructions.txt",
				MimeType: "text/plain; charset=utf-8",
				Parents:  []string{"child_test_folder"},
			},
		},
	}

	for _, c := range cases {
		fileData, err := ValidateFile(c.path)
		if err != nil {
			t.Errorf("Fail: unable to validate file: %v\n", err)
		}

		filePointer := writeUploadMetadata(fileData, c.targetFolder)

		if filePointer.Name != c.expected.Name {
			t.Errorf(
				"Fail: drive.File.Name does not match expected:\nExpected: %v\nActual: %v\n",
				c.expected.Name,
				filePointer.Name,
			)
		}

		if filePointer.MimeType != c.expected.MimeType {
			t.Errorf(
				"Fail: drive.File.MimeType does not match expected:\nExpected: %v\nActual: %v\n",
				c.expected.MimeType,
				filePointer.MimeType,
			)
		}

		t.Logf("Expected parents: %v\nActual parents: %v", c.expected.Parents, filePointer.Parents)
		if c.targetFolder != "" && filePointer.Parents[0] != c.expected.Parents[0] {
			t.Errorf(
				"Fail: drive.File.Parents does not match expected:\nExpected: %v\nActual: %v\n",
				c.expected.Parents[0],
				filePointer.Parents[0],
			)
		}
	}
}

func TestUploadFile(t *testing.T) {
	var config apiConfig
	err := config.setClient()
	if err != nil {
		log.Fatalf("Error setting client: %v\n", err)
	}

	err = config.setService()
	if err != nil {
		log.Fatalf("Error setting service: %v\n", err)
	}

	cases := []struct {
		path   string
		target string
	}{
		{
			path:   "./test_resources/bubbletea_gui.gif",
			target: "test_folder",
		}, {
			path:   "./test_resources/goose_sqlc-instructions.txt",
			target: "/test_folder/child_test_folder",
		},
	}

	for _, c := range cases {
		_, err = config.uploadFile(c.path, c.target)
		if err != nil {
			t.Errorf("Fail: error creating file on drive: %v", err)
		}
	}
}
