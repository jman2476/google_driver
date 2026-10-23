package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"testing"
	"time"
)

func genTimeStampFilename(base, ext string) string {
	currTime := time.Now()
	return fmt.Sprintf("%s-%v:%v%s", base, currTime.Hour(), currTime.Minute(), ext)
}

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
		path      string
		target    string
		basename  string
		extension string
		in        io.Reader
	}{
		{
			path:      "./test_resources/bubbletea_gui_for_delete.gif",
			target:    "test_folder",
			basename:  "bubbletea_gui_for_delete",
			extension: "gif",
			in:        strings.NewReader("y\ny\n"),
		}, {
			path:      "./test_resources/goose_sqlc-instructions_for_delete.txt",
			target:    "/test_folder/child_test_folder",
			basename:  "goose_sqlc-instructions_for_delete.txt",
			extension: "txt",
			in:        strings.NewReader("y\ny\n"),
		},
	}

	for _, c := range cases {
		config.in = c.in
		_, err := config.uploadFile(c.path, c.target)
		if err != nil {
			t.Errorf("Fail case %s: error creating file on drive: %v", c.path, err)
		}

		_, err = config.removeFile(
			genTimeStampFilename(c.basename, c.extension), c.target)
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
		path        string
		target      string
		filename    string
		expectedErr error
		in          io.Reader
	}{
		{
			path:        "./test_resources/bubbletea_gui.gif",
			target:      "test_folder",
			filename:    "bubbletea_gui.gif",
			expectedErr: ErrTooManyFilesFound,
			in:          strings.NewReader("y\ny\n"),
		}, {
			path:        "./test_resources/goose_sqlc-instructions_for_delete.txt",
			target:      "/test_folder/child_test_folder",
			filename:    "goose_for_delete.txt",
			expectedErr: ErrFileNotFound,
			in:          strings.NewReader("y\ny\n"),
		},
	}

	for _, c := range cases {
		config.in = c.in
		_, err := config.uploadFile(c.path, c.target)
		if err != nil {
			t.Errorf("Fail case %s: error creating file on drive: %v", c.path, err)
		}

		_, err = config.removeFile(c.filename, c.target)
		if err.Error() != fmt.Sprintf("%v: %s", c.expectedErr, c.filename) {
			t.Errorf(
				"Fail case %s:\nExpected: %v\nActual: %v",
				c.path, c.expectedErr, err,
			)
		}
	}

	_, err = config.removeFile(
		"goose_sqlc-instructions_for_delete.txt",
		"/test_folder/child_test_folder",
	)
	if err != nil {
		t.Errorf("Fail: can't clean up test cases: %v", err)
	}
}
