package util

import (
	"os"
	"testing"

	"fmt"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAbsPathify(t *testing.T) {
	const (
		ENV_VAR   = "HIDDEN_VAR"
		ENV_VALUE = "/foo"
	)
	Convey(fmt.Sprintf(
		`Given
 - the environment variable %s = %s
 - the current directory '%s'`, ENV_VAR, ENV_VALUE, mustGetwd()), t, func() {
		os.Setenv(ENV_VAR, ENV_VALUE)
		wd := mustGetwd()

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

		for _, tt := range data {
			Convey(fmt.Sprintf("When when we get the absolute path of '%s'", tt.path), func() {
				result := AbsPathify(tt.path)

				Convey(fmt.Sprintf("The absolute path should be '%s'", tt.expectedResult), func() {
					So(result, ShouldEqual, tt.expectedResult)
				})
			})
		}
	})
}

func TestExists(t *testing.T) {
	Convey(fmt.Sprintf(
		`Given
 - the current directory '%s'`, mustGetwd()), t, func() {
		wd, _ := os.Getwd()
		var data = []struct {
			path           string // path to test
			expectedResult bool   // expected result
		}{
			{wd, true},
			{".", true},
			{"/non/existing/path", false},
		}

		for _, tt := range data {
			Convey(fmt.Sprintf("When when we test '%s' for existence", tt.path), func() {
				result := Exists(tt.path)

				Convey(fmt.Sprintf("The result should be '%t'", tt.expectedResult), func() {
					So(result, ShouldEqual, tt.expectedResult)
				})
			})
		}
	})
}

func mustGetwd() string {
	wd, _ := os.Getwd()

	return wd
}
