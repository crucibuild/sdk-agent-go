// Copyright (C) 2016 Christophe Camel, Jonathan Pigrée
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"os"
	"testing"

	"fmt"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAbsPathify(t *testing.T) {
	const (
		EnvVar   = "HIDDEN_VAR"
		EnvValue = "/foo"
	)
	Convey(fmt.Sprintf(
		`Given
 - the environment variable %s = %s
 - the current directory '%s'`, EnvVar, EnvValue, mustGetwd()), t, func() {
		os.Setenv(EnvVar, EnvValue)
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
