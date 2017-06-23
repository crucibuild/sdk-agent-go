package util

import (
	"os"
	"testing"
)

func TestAbsPathify(t *testing.T) {
	// Arrange
	os.Setenv("HIDDEN_VAR", "/foo")
	wd, _ := os.Getwd()
	var data = []struct {
		path           string // path to test
		expectedResult string // expected absolute path
	}{
		{"", wd},
		{"/foo/bar", "/foo/bar"},
		{"/foo//bar", "/foo/bar"},
		{"/foo/../bar", "/bar"},
		{"foo/bar", wd + "/foo/bar"},
		{"$HOME", UserHomeDir()},
		{"$HOME/foo/bar", UserHomeDir() + "/foo/bar"},
		{"$HIDDEN_VAR/bar", "/foo/bar"},
	}

	for i, tt := range data {
		// Act
		result := AbsPathify(tt.path)

		// Assert
		if result != tt.expectedResult {
			t.Errorf("test %d: Bad absolute path returned. Expected '%s', got '%s'", i, tt.expectedResult, result)
		}
	}
}

func TestExists(t *testing.T) {
	// Arrange
	wd, _ := os.Getwd()
	var data = []struct {
		path           string // path to test
		expectedResult bool   // expected result
	}{
		{wd, true},
		{".", true},
		{"/non/existing/path", false},
	}

	for i, tt := range data {
		// Act
		result := Exists(tt.path)

		// Assert
		if result != tt.expectedResult {
			t.Errorf("test %d: Expected '%t', got '%t' for path '%s'", i, tt.expectedResult, result, tt.path)
		}
	}
}
