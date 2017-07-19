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
	"encoding/json"
	"fmt"
	"github.com/crucibuild/sdk-agent-go/agentiface"
	"github.com/crucibuild/sdk-agent-go/util"
)

// MimeTypeJSONSchema represents the mime type we use for JSON schemas.
const MimeTypeJSONSchema = "application/schema+json"

const jsonID = "id"
const jsonTitle = "title"

// JSONSchema represents a JSON Schema with all its metadata.
type JSONSchema struct {
	id    string
	title string
	raw   string
}

// ID returns the JSONSchema ID.
func (s *JSONSchema) ID() string {
	return s.id
}

// Title returns the JSONSchema title.
func (s *JSONSchema) Title() string {
	return s.title
}

// MimeType returns the JSONSchema mime type.
func (*JSONSchema) MimeType() string {
	return MimeTypeJSONSchema
}

// Raw returns the raw JSONSchema.
func (s *JSONSchema) Raw() string {
	return s.raw
}

// Decode unserializes data using the JSONSchema registered.
func (s *JSONSchema) Decode(o []byte, t agentiface.Type) (interface{}, error) {
	// Create a new record to decode data into
	decodedRecord := util.New(t.Type())

	// decode
	err := json.Unmarshal(o, decodedRecord)

	return decodedRecord, err
}

// Code serializes data using the JSONSchema registered.
func (s *JSONSchema) Code(o interface{}) ([]byte, error) {
	return json.Marshal(o)
}

// LoadJSONSchema loads the given Json Schema and returns a schema instance
func LoadJSONSchema(rawSchema string) (agentiface.Schema, error) {
	// The given json schema is a json, so load it
	var decoded interface{}
	err := json.Unmarshal([]byte(rawSchema), &decoded)

	if err != nil {
		return nil, err
	}

	schema := decoded.(map[string]interface{})

	id, ok := schema[jsonID]

	if !ok {
		return nil, fmt.Errorf("id (key: %s) was expected in schema", jsonID)
	}

	switch id.(type) {
	case string:
		// ok
		break
	default:
		return nil, fmt.Errorf("id (key: %s) value must be a JSON string in schema", jsonID)
	}

	return &JSONSchema{
		id:    id.(string),
		title: schema[jsonTitle].(string),
		raw:   rawSchema,
	}, nil
}
