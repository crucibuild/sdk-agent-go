/* Copyright (C) 2016 Christophe Camel, Jonathan Pigrée
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type aStruct struct {
	Value string
}

func TestGetStructType(t *testing.T) {
	// Arrange
	var data = []struct {
		instance     interface{}  // instance to test
		expectedType reflect.Type // expected type
		expectError  bool
	}{
		{aStruct{}, reflect.TypeOf(aStruct{}), false},
		{&aStruct{}, reflect.TypeOf(aStruct{}), false},
		{"foo", nil, true},
	}

	for i, tt := range data {
		// Act
		result, err := GetStructType(tt.instance)

		// Assert
		assert.Equal(t, tt.expectError, err != nil,
			"test %d: Expected error and returned error must match", i)

		assert.Equal(t, tt.expectedType, result,
			"test %d: Type must match", i)
	}
}

func TestNew(t *testing.T) {
	// Arrange
	type ComplexStruct struct {
		V1 string
		V2 map[int]string
		V3 []int
		V4 chan int
		V5 aStruct
		V6 *aStruct
	}

	tpe := reflect.TypeOf(ComplexStruct{})

	// Act
	v := New(tpe)

	// Assert
	assert.NotNil(t, v, "Value must not be nil")

	// test type conversion
	aStruct := v.(*ComplexStruct)

	// fields must be initialized
	assert.NotNil(t, aStruct.V2, "V2 must be initialized")

	assert.NotNil(t, aStruct.V3, "V3 must be initialized")

	assert.NotNil(t, aStruct.V4, "V4 must be initialized")

	assert.NotNil(t, aStruct.V6, "V6 must be initialized")
}
