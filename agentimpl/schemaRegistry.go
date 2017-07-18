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
)

// SchemaRegistry represents a registry for schemas.
type SchemaRegistry struct {
	schemas map[string]agentiface.Schema
}

// NewSchemaRegistry creates a new instance of SchemaRegistry.
// nolint: unparam, parameter a is reserved for a future usage
func NewSchemaRegistry(a agentiface.Agent) *SchemaRegistry {
	return &SchemaRegistry{
		schemas: make(map[string]agentiface.Schema),
	}
}

// SchemaRegister registers a schema in the registry.
func (s *SchemaRegistry) SchemaRegister(schema agentiface.Schema) (string, error) {
	s.schemas[schema.ID()] = schema

	return schema.ID(), nil
}

// SchemaGetByID returns a schema which id matches the one in parameter.
func (s *SchemaRegistry) SchemaGetByID(id string) (agentiface.Schema, error) {
	schema, ok := s.schemas[id]

	if !ok {
		return nil, fmt.Errorf("No schema found in the registry with id '%s'", id)
	}

	return schema, nil
}

// SchemaListIds returns a map of <id, schema> known by the registry.
func (s *SchemaRegistry) SchemaListIds() []string {
	values := make([]string, len(s.schemas))

	i := 0
	for k := range s.schemas {
		values[i] = k
		i++
	}

	return values
}

// SchemaUnregister remove a schema from the registry.
func (s *SchemaRegistry) SchemaUnregister(id string) error {
	if !s.SchemaExist(id) {
		return fmt.Errorf("No schema found in the registry with id '%s'", id)
	}

	delete(s.schemas, id)

	return nil
}

// SchemaExist returns true if the key match a schema known by the registry.
func (s *SchemaRegistry) SchemaExist(key string) bool {
	_, ok := s.schemas[key]

	return ok
}
