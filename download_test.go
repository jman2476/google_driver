package main

import "testing"

func TestDownloadFile(t *testing.T) {
	cases := []struct {
		input struct {
			path   string
			target string
			flags  map[string]string
		}
		expected    string
		expectedErr error
	}{}
}
