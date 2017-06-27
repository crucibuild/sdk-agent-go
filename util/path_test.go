package util

import (
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, tt.expectedResult, result,
			"test %d: Absolute path must match (for path '%s')", i, tt.path)
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
		assert.Equal(t, tt.expectedResult, result,
			"test %d: Path existence must match (for path '%s')", i, tt.path)

	}
}
