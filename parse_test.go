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
			input:    "",
			expected: arguments{},
		}, {
			input:    "",
			expected: arguments{},
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
