package main

import "testing"

func TestRemoveFile(t *testing.T) {
	// Test steps:
	//  1-Add file xyz.md
	//  2-Delete file xyz.md
	//  3-Check for file xyz.md
}

func TestRemoveFileErrors(t *testing.T) {
	// Test steps:
	//  1-Add file xyz.md twice
	//  2-Attempt delete xyz.md
	//  3-Check error
	//  4-Attempt delete non-existent file
	//  5-Check error
}
