package main

import (
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
			targetFolder: "/Testing/",
			expected: &drive.File{
				Name:     "bubbletea_gui.gif",
				MimeType: "image/gif",
				Parents:  []string{"/Testing/"},
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
			targetFolder: "",
			expected: &drive.File{
				Name:     "goose_sqlc-instructions.txt",
				MimeType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, c := range cases {
		fileData, err := ValidateFile(c.path)
		if err != nil {
			t.Errorf("Fail: unable to validate file: %v", err)
		}

		filePointer := WriteUploadMetadata(fileData, c.targetFolder)

		if filePointer.Name != c.expected.Name {
			t.Errorf(
				"Fail: drive.File.Name does not match expected:\nExpected: %v\nActual: %v",
				c.expected.Name,
				filePointer.Name,
			)
		}

		if filePointer.MimeType != c.expected.MimeType {
			t.Errorf(
				"Fail: drive.File.Name does not match expected:\nExpected: %v\nActual: %v",
				c.expected.MimeType,
				filePointer.MimeType,
			)
		}

		if c.targetFolder != "" && filePointer.Parents[0] != c.expected.Parents[0] {
			t.Errorf(
				"Fail: drive.File.Name does not match expected:\nExpected: %v\nActual: %v",
				c.expected.Parents[0],
				filePointer.Parents[0],
			)
		}
	}
}
