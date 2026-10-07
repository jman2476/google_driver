package main

import (
	"fmt"
	"log"
	"testing"
)

func TestParsePath(t *testing.T) {
	cases := []struct {
		input       string
		expected    []DriveFolder
		expectedErr error
	}{
		{
			input: "first/second/third",
			expected: []DriveFolder{
				{Name: "first", Parent: ""},
				{Name: "second", Parent: "first"},
				{Name: "third", Parent: "second"},
			},
			expectedErr: nil,
		}, {
			input: "/first/second/third/",
			expected: []DriveFolder{
				{Name: "first", Parent: ""},
				{Name: "second", Parent: "first"},
				{Name: "third", Parent: "second"},
			},
			expectedErr: nil,
		}, {
			input:       "//first/",
			expected:    []DriveFolder{},
			expectedErr: ErrParseDirs,
		}, {
			input:       "/",
			expected:    []DriveFolder{},
			expectedErr: ErrEmptyDir,
		}, {
			input:       "//",
			expected:    []DriveFolder{},
			expectedErr: ErrParseDirs,
		}, {
			input:       "/first//",
			expected:    []DriveFolder{},
			expectedErr: ErrParseDirs,
		},
	}

	for _, c := range cases {
		result, err := ParsePath(c.input)

		if len(result) != len(c.expected) {
			t.Errorf(
				"Fail: Parse path result length doesn't match\nCase: %v\nExpected: %v\nActual: %v",
				c.input, len(c.expected), len(result),
			)
		}

		for i, r := range result {
			if r != c.expected[i] {
				t.Errorf(
					"Fail: ParsePath result doesn't match expected\nCase: %v\nExpected: %v\nActual: %v",
					c.input, c.expected[i], r,
				)
			}
		}

		if err != c.expectedErr {
			t.Errorf(
				"Fail: ParsePath error doesn't match expected\nCase: %v\nExpected: %v\nActual: %v",
				c.input, c.expectedErr, err,
			)
		}
	}
}

func TestFindFolder(t *testing.T) {
	var config apiConfig
	err := config.setClient()
	if err != nil {
		log.Fatalf("Error setting client: %v\n", err)
	}

	err = config.setService()
	if err != nil {
		log.Fatalf("Error setting service: %v\n", err)
	}

	cases := "test_folder"

	folder, err := config.FindFolder(cases)
	if err != nil {
		t.Logf("Fail: TFF error: %v", err)
	}

	fmt.Printf("Folder: %s", folder.Name)
}

func TestCreateFolder(t *testing.T) {
	var config apiConfig
	err := config.setClient()
	if err != nil {
		t.Errorf("error setting client: %v\n", err)
	}

	err = config.setService()
	if err != nil {
		t.Errorf("error setting service: %v\n", err)
	}

	path := "test_folder/child_test_folder"
	folder, err := config.CreateFolder(path)
	if err != nil {
		t.Errorf("Fail: error creating folder on drive: %v", err)
	}

	t.Logf(
		"Created folder %v\nParent: %v\nID: %v\n",
		folder.Name, folder.Parents, folder.Id,
	)
}

func TestDeleteFolder(t *testing.T) {

}
