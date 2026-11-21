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
				tool:       "mv",
				positional: []string{"./test_resources/chess.mov", "/test_folder"},
			},
			expectedErr: nil,
		}, {
			input: "google_driver rm /test_folder/chess.mov -y -y",
			expected: arguments{
				tool:       "rm",
				positional: []string{"/test_folder/chess.mov"},
				flags:      map[string]string{"y": "true"},
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

		if len(parsed.positional) != len(c.expected.positional) {
			t.Errorf("Fail: incorrect number of positional arguments\nCase: %s\nExpected: %d\nActual: %d\n",
				c.input,
				len(c.expected.positional),
				len(parsed.positional),
			)
		} else {
			for i, p := range parsed.positional {
				if p != c.expected.positional[i] {
					t.Errorf(
						"Fail: positional arguments don't match expected\nCase: %s\nExpected: %v\nActual: %v\n",
						c.input,
						c.expected.positional,
						parsed.positional,
					)
					break
				}
			}
		}

		if parsed.tool != c.expected.tool {
			t.Errorf(
				"Fail: tool does not match expected\nCase: [%s]\nExpected: %v\nActual: %v\n",
				c.input,
				c.expected.tool,
				parsed.tool,
			)
		}

		if len(parsed.flags) != len(c.expected.flags) {
			t.Errorf(
				"Fail: incorrect number of flags\nCase: %s\nExpected: %d\nActual: %d\n",
				c.input,
				len(c.expected.flags),
				len(parsed.flags),
			)
		} else {
			for key, val := range parsed.flags {
				eVal, ok := c.expected.flags[key]
				if !ok {
					t.Errorf(
						"Fail: unexpected flag found: %s\nCase: %s\nExpected: none",
						val, c.input,
					)
					break
				}
				if val != eVal {
					t.Errorf(
						"Fail: value for flag %s not consistent\nCase: %s\nExpected: %s\nActual: %s\n",
						key, c.input,
						eVal, val,
					)
					break
				}
			}
		}
	}
}
