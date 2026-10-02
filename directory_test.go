package main

import "testing"

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

}

func TestCreateFolder(t *testing.T) {

}

func TestDeleteFolder(t *testing.T) {

}
