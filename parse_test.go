package main

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		input       string
		expected    arguments
		expectedErr error
	}{
		{
			input: "google_driver mv ./test_resources/chess.mov /test_folder",
			expected: arguments{
				tool:   "mv",
				source: "./test_resources/chess.mov",
				target: "/test_folder/",
			},
			expectedErr: nil,
		}, {
			input: "google_driver rm /test_folder/chess.mov -y -y",
			expected: arguments{
				tool:   "rm",
				source: "/test_folder/chess.mov",
				target: "/test_folder",
			},
		},
	}

	for _, c := range cases {
		args := strings.Split(c.input, " ")

		parsed, err := parseArgs(args)
		if err != nil && err != c.expectedErr {
			t.Errorf(
				"Fail: error does not match expected\nCase: [%s]\nExpected: %v\nActual: %v\n",
				c.input,
				c.expectedErr,
				err,
			)
		}

		if parsed.source != c.expected.source {
			t.Errorf(
				"Fail: source does not match expected\nCase: [%s]\nExpected: %v\nActual: %v\n",
				c.input,
				c.expected.source,
				parsed.source,
			)
		}

		if parsed.tool != c.expected.tool {
			t.Errorf(
				"Fail: tool does not match expected\nCase: [%s]\nExpected: %v\nActual: %v\n",
				c.input,
				c.expected.tool,
				parsed.tool,
			)
		}
		if parsed.target != c.expected.target {
			t.Errorf(
				"Fail: target does not match expected\nCase: [%s]\nExpected: %v\nActual: %v\n",
				c.input,
				c.expected.target,
				parsed.target,
			)
		}
	}
}
