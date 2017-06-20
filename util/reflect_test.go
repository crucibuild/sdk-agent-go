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
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	// Arrange
	type AStruct struct {
		Value string
	}

	x := AStruct{}
	tpe, err := GetStructType(x)

	if err != nil {
		t.Error(err.Error())
	}

	// Act
	v := New(tpe)

	// Assert
	if v == nil {
		t.Error("Value is nil")
	}

	// test type conversion
	aStruct := v.(*AStruct)
	aStruct.Value = "test"

	// test field
	ntpe, err := GetStructType(v)
	if err != nil {
		t.Error(err.Error())
	}

	if ntpe.NumField() != 1 {
		t.Error("Incorrect number of fields")
	}

	f := ntpe.Field(0)
	if f.Name != "Value" {
		t.Error("Incorrect field name")
	}

	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		t.Error("Not applicable for non-pointer types or nil")
	}

}
