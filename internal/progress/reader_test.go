package progress

import "testing"

func TestFormatRate(t *testing.T) {
	cases := []struct {
		input    float64
		expected string
	}{
		{
			input:    5.24,
			expected: "5.2 B/s",
		},
		{
			input:    2000.0,
			expected: "2.0 KB/s",
		},
		{
			input:    432334432.034,
			expected: "432.3 MB/s",
		},
		{
			input:    4473382848.0,
			expected: "4.5 GB/s",
		},
		{
			input:    56645698744585,
			expected: "56.6 TB/s",
		},
		{
			input:    432.29,
			expected: "432.3 B/s",
		},
	}

	for _, c := range cases {
		result := formatRate(c.input)

		if result != c.expected {
			t.Errorf("Fail: Case %v\nExpected: %v\nActual: %v\n", c.input, c.expected, result)
		}
	}
}
