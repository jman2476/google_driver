package main

import (
	"log"
	"testing"
)

func TestRemoveFile(t *testing.T) {
	// Test steps:
	//  1-Add file xyz.md
	//  2-Delete file xyz.md
	//  3-Check for file xyz.md
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
		path     string
		target   string
		filename string
	}{
		{
			path:     "./test_resources/bubbletea_gui_for_delete.gif",
			target:   "test_folder",
			filename: "bubbletea_gui_for_delete.gif",
		}, {
			path:     "./test_resources/goose_sqlc-instructions.txt",
			target:   "/test_folder/child_test_folder",
			filename: "goose_sqlc-instructions.txt",
		},
	}

	for _, c := range cases {
		_, err := config.uploadFile(c.path, c.target)
		if err != nil {
			t.Errorf("Fail case %s: error creating file on drive: %v", c.path, err)
		}

		_, err = config.removeFile(c.filename, c.target)
		if err != nil {
			t.Errorf("Fail case %s: error deleting file on drive: %v", c.path, err)
		}
	}
}

func TestRemoveFileErrors(t *testing.T) {
	// Test steps:
	//  1-Add file xyz.md twice
	//  2-Attempt delete xyz.md
	//  3-Check error
	//  4-Attempt delete non-existent file
	//  5-Check error
}
