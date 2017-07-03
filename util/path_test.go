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
