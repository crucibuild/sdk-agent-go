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
	"github.com/pkg/errors"
)

type typeStruct struct {
	name string
	t    reflect.Type
}

type TypeRegistry struct {
	typesByName map[string]agentiface.Type
	typesByType map[reflect.Type]agentiface.Type
}

func NewTypeRegistry(a agentiface.Agent) *TypeRegistry {
	return &TypeRegistry{
		typesByName: make(map[string]agentiface.Type),
		typesByType: make(map[reflect.Type]agentiface.Type),
	}
}

func (t *typeStruct) Name() string {
	return t.name
}

func (t *typeStruct) Type() reflect.Type {
	return t.t
}

func NewTypeFromInterface(name string, i interface{}) (agentiface.Type, error) {
	t, err := util.GetStructType(i)

	if err != nil {
		return nil, err
	}

	return NewTypeFromType(name, t), nil
}

func NewTypeFromType(name string, t reflect.Type) agentiface.Type {
	return &typeStruct{
		name: name,
		t:    t,
	}
}

func (s *TypeRegistry) TypeRegister(t agentiface.Type) (string, error) {
	s.typesByName[t.Name()] = t
	s.typesByType[t.Type()] = t

	return t.Name(), nil
}

func (s *TypeRegistry) TypeGetByName(key string) (agentiface.Type, error) {
	t, ok := s.typesByName[key]

	if !ok {
		return nil, errors.New(fmt.Sprintf("No type found in the registry with key '%s'", key))
	}

	return t, nil
}

func (s *TypeRegistry) TypeGetByType(v reflect.Type) (agentiface.Type, error) {
	t, ok := s.typesByType[v]

	if !ok {
		return nil, errors.New(fmt.Sprintf("No type found in the registry with type '%s'", v.Name()))
	}

	return t, nil
}

func (s *TypeRegistry) TypeListNames() []string {
	values := make([]string, len(s.typesByName))

	i := 0
	for k := range s.typesByName {
		values[i] = k
		i++
	}

	return values
}

func (s *TypeRegistry) TypeUnregister(key string) error {
	t, err := s.TypeGetByName(key)

	if err != nil {
		return err
	}

	delete(s.typesByName, key)
	delete(s.typesByType, t.Type())

	return nil
}

func (s *TypeRegistry) TypeExist(key string) bool {
	_, ok := s.typesByName[key]

	return ok
}
