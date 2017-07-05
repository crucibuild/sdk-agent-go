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

package agentimpl

import (
	"fmt"
	"github.com/crucibuild/sdk-agent-go/agentiface"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

func TestNewTypeFromInterface(t *testing.T) {
	const (
		STRUCT_NAME = "foo"
	)
	Convey(fmt.Sprintf("Given a struct instance named '%s'", STRUCT_NAME), t, func() {
		i := struct {
			Value string
		}{
			"bar",
		}
		expectedType := reflect.TypeOf(i)

		Convey(fmt.Sprintf("When when we call the NewTypeFromInterface() function"), func() {
			tpe, err := NewTypeFromInterface(STRUCT_NAME, i)

			Convey("No error should occur", func() {
				So(err, ShouldBeNil)
			})

			Convey("Type should not be nil", func() {
				So(tpe, ShouldNotBeNil)
			})

			Convey(fmt.Sprintf("Name of the type should be equal to '%s'", STRUCT_NAME), func() {
				So(tpe.Name(), ShouldEqual, STRUCT_NAME)
			})

			Convey(fmt.Sprintf("Type should be equal to '%s'", expectedType.Kind()), func() {
				So(tpe.Type(), ShouldEqual, expectedType)
			})
		})
	})
}

func TestNewTypeFromType(t *testing.T) {
	var (
		TYPE_NAME = "foo"
		TYPE      = reflect.TypeOf("")
	)
	Convey(fmt.Sprintf("Given a type '%s' named '%s'", TYPE.Kind(), TYPE_NAME), t, func() {
		Convey(fmt.Sprintf("When when we call the NewTypeFromType() function"), func() {
			tpe := NewTypeFromType(TYPE_NAME, TYPE)

			Convey("Type should not be nil", func() {
				So(tpe, ShouldNotBeNil)
			})

			Convey(fmt.Sprintf("Name of the type should be equal to '%s'", TYPE_NAME), func() {
				So(tpe.Name(), ShouldEqual, TYPE_NAME)
			})

			Convey(fmt.Sprintf("Type should be equal to '%s'", TYPE.Kind()), func() {
				So(tpe.Type(), ShouldEqual, TYPE)
			})
		})
	})
}

func TestNewTypeRegistry(t *testing.T) {
	Convey("Given an agent", t, func() {
		var agent agentiface.Agent = nil // not used

		Convey(fmt.Sprintf("When when we create a new type registry"), func() {
			registry := NewTypeRegistry(agent)

			Convey("The registry should no be nil", func() {
				So(registry, ShouldNotBeNil)
			})
			Convey("The registry should no be empty", func() {
				So(len(registry.TypeListNames()), ShouldEqual, 0)
			})
		})
	})
}

func TestRegisterANewType(t *testing.T) {
	Convey("Given an empty registry", t, func() {
		var agent agentiface.Agent = nil // not used
		expectedType := NewTypeFromType("foo", reflect.TypeOf(""))
		registry := NewTypeRegistry(agent)

		Convey(fmt.Sprintf("When when we register the new type '%s'", expectedType.Name()), func() {

			registry.TypeRegister(expectedType)

			Convey("The size of the registry should be equal to 1", func() {
				So(len(registry.TypeListNames()), ShouldEqual, 1)
			})

			Convey(fmt.Sprintf("The registry should contain the type '%s'", expectedType.Name()), func() {
				So(registry.TypeExist(expectedType.Name()), ShouldBeTrue)
			})

			Convey(fmt.Sprintf("When the type '%s' is retrieved by its name in the registry", expectedType.Name()), func() {
				tpe, err := registry.TypeGetByName("foo")

				Convey("No error should occur", func() {
					So(err, ShouldBeNil)
				})
				Convey(fmt.Sprintf("Type retrieved should be '%s'", expectedType.Name()), func() {
					So(tpe, ShouldEqual, expectedType)
				})
			})

			Convey(fmt.Sprintf("When the type '%s' is retrieved by its type in the registry", expectedType.Name()), func() {
				tpe, err := registry.TypeGetByType(expectedType.Type())

				Convey("No error should occur", func() {
					So(err, ShouldBeNil)
				})
				Convey(fmt.Sprintf("Type retrieved should be '%s'", expectedType.Name()), func() {
					So(tpe, ShouldEqual, expectedType)
				})
			})
		})
	})
}

func TestUnregisterAType(t *testing.T) {
	var (
		TYPE_NAME = "foo"
		TYPE      = reflect.TypeOf("")
	)
	Convey(fmt.Sprintf(`Given a registry containing only the type '%s' (%s)`, TYPE_NAME, TYPE.Kind()), t, func() {
		var agent agentiface.Agent = nil
		expectedType := NewTypeFromType(TYPE_NAME, TYPE)
		registry := NewTypeRegistry(agent)
		registry.TypeRegister(expectedType)

		Convey(fmt.Sprintf("When when we unregister the type '%s'", expectedType.Name()), func() {
			err := registry.TypeUnregister("foo")

			Convey("No error should occur", func() {
				So(err, ShouldBeNil)
			})
			Convey("The registry should be empty", func() {
				So(len(registry.TypeListNames()), ShouldEqual, 0)
			})
		})
	})
}
