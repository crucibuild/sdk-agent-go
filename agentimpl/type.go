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
	"reflect"

	"github.com/crucibuild/sdk-agent-go/agentiface"
	"github.com/crucibuild/sdk-agent-go/util"
)

type typeStruct struct {
	name string
	t    reflect.Type
}

// TypeRegistry is a registry referencing types.
type TypeRegistry struct {
	typesByName map[string]agentiface.Type
	typesByType map[reflect.Type]agentiface.Type
}

// NewTypeRegistry returns a new instance of TypeRegistry from an agent.
func NewTypeRegistry(a agentiface.Agent) *TypeRegistry {
	return &TypeRegistry{
		typesByName: make(map[string]agentiface.Type),
		typesByType: make(map[reflect.Type]agentiface.Type),
	}
}

// Name returns the name of a Type.
func (t *typeStruct) Name() string {
	return t.name
}

// Name returns a reflect Type of a Type.
func (t *typeStruct) Type() reflect.Type {
	return t.t
}

// NewTypeFromInterface returns a new instance of Type from a name and an interface.
func NewTypeFromInterface(name string, i interface{}) (agentiface.Type, error) {
	t, err := util.GetStructType(i)

	if err != nil {
		return nil, err
	}

	return NewTypeFromType(name, t), nil
}

// NewTypeFromType returns a Type from a name and a Type reflection.
func NewTypeFromType(name string, t reflect.Type) agentiface.Type {
	return &typeStruct{
		name: name,
		t:    t,
	}
}

// TypeRegister registers a Type into the registry.
func (s *TypeRegistry) TypeRegister(t agentiface.Type) (string, error) {
	s.typesByName[t.Name()] = t
	s.typesByType[t.Type()] = t

	return t.Name(), nil
}

// TypeGetByName returns a Type which name's match key.
func (s *TypeRegistry) TypeGetByName(key string) (agentiface.Type, error) {
	t, ok := s.typesByName[key]

	if !ok {
		return nil, fmt.Errorf("No type found in the registry with key '%s'", key)
	}

	return t, nil
}

// TypeGetByType returns a Type from a reflect.Type.
func (s *TypeRegistry) TypeGetByType(v reflect.Type) (agentiface.Type, error) {
	t, ok := s.typesByType[v]

	if !ok {
		return nil, fmt.Errorf("No type found in the registry with type '%s'", v.Name())
	}

	return t, nil
}

// TypeListNames returns a map(<id, name>) of all the registered types.
func (s *TypeRegistry) TypeListNames() []string {
	values := make([]string, len(s.typesByName))

	i := 0
	for k := range s.typesByName {
		values[i] = k
		i++
	}

	return values
}

// TypeUnregister unregister the Type key from the registry.
func (s *TypeRegistry) TypeUnregister(key string) error {
	t, err := s.TypeGetByName(key)

	if err != nil {
		return err
	}

	delete(s.typesByName, key)
	delete(s.typesByType, t.Type())

	return nil
}

// TypeExist returns true if a type matching key is known to the registry.
func (s *TypeRegistry) TypeExist(key string) bool {
	_, ok := s.typesByName[key]

	return ok
}
