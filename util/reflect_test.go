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

	"fmt"
	. "github.com/smartystreets/goconvey/convey"
)

type aStruct struct {
	Value string
}

func TestGetStructType(t *testing.T) {
	Convey("Given a set of interfaces", t, func() {
		var data = []struct {
			name         string
			instance     interface{}  // instance to test
			expectedType reflect.Type // expected type
			expectError  bool
		}{
			{"a struct", aStruct{}, reflect.TypeOf(aStruct{}), false},
			{"a pointer to a struct", &aStruct{}, reflect.TypeOf(aStruct{}), false},
			{"a string", "foo", nil, true},
		}

		for _, tt := range data {
			Convey(fmt.Sprintf("When when we get the type from %s", tt.name), func() {
				result, err := GetStructType(tt.instance)

				if tt.expectError {
					Convey("An error should occur", func() {
						So(err, ShouldNotBeNil)
					})

					Convey("Value should be nil", func() {
						So(result, ShouldBeNil)
					})
				} else {
					Convey("No error should occur", func() {
						So(err, ShouldBeNil)
					})

					Convey(fmt.Sprintf("Value should be '%s'", tt.expectedType.Name()), func() {
						So(result, ShouldEqual, tt.expectedType)
					})
				}
			})
		}
	})
}

func TestNew(t *testing.T) {
	Convey("Given a struct", t, func() {
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

		Convey("When when we call the New() function", func() {
			v := New(tpe)

			Convey("Created instance should not be nil", func() {
				So(v, ShouldNotBeNil)
			})

			Convey(fmt.Sprintf("Created instance should have the type '%s'", tpe.Name()), func() {
				So(v, ShouldHaveSameTypeAs, &ComplexStruct{})
			})

			Convey("Fields of the created instance should be initialized", func() {
				s := v.(*ComplexStruct)

				So(s.V2, ShouldNotBeNil)
				So(s.V3, ShouldNotBeNil)
				So(s.V4, ShouldNotBeNil)
				So(s.V6, ShouldNotBeNil)
			})
		})
	})
}
