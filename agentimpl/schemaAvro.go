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
	"bytes"
	"github.com/crucibuild/sdk-agent-go/agentiface"
	"github.com/crucibuild/sdk-agent-go/util"
	"github.com/elodina/go-avro"
)

// MimeTypeAvroSchema represents the mime type we use for AVRO schemas.
const MimeTypeAvroSchema = "application/js+avro"

// AvroSchema represents an AVRO schema with all its metadata.
type AvroSchema struct {
	id     string
	title  string
	raw    string
	schema avro.Schema
}

// ID returns the AvroSchema ID.
func (s *AvroSchema) ID() string {
	return s.id
}

// Title returns the AvroSchema title.
func (s *AvroSchema) Title() string {
	return s.title
}

// MimeType returns the AvroSchema mime type.
func (*AvroSchema) MimeType() string {
	return MimeTypeAvroSchema
}

// Raw returns the raw AvroSchema.
func (s *AvroSchema) Raw() string {
	return s.raw
}

// Decode unserializes data using the AvroSchema registered.
func (s *AvroSchema) Decode(o []byte, t agentiface.Type) (interface{}, error) {
	// Create a new Decoder with the data
	decoder := avro.NewBinaryDecoder(o)

	// Create a new record to decode data into
	decodedRecord := util.New(t.Type())

	// Read data into a given record with a given Decoder.
	reader := avro.NewSpecificDatumReader()
	reader.SetSchema(s.schema)

	// decode
	err := reader.Read(decodedRecord, decoder)

	return decodedRecord, err
}

// Code serializes data using the AvroSchema registered.
func (s *AvroSchema) Code(o interface{}) ([]byte, error) {
	// encode command
	writer := avro.NewSpecificDatumWriter()
	writer.SetSchema(s.schema)
	buffer := new(bytes.Buffer)
	encoder := avro.NewBinaryEncoder(buffer)

	err := writer.Write(o, encoder)

	return buffer.Bytes(), err
}

// LoadAvroSchema loads the given raw Avro definition and returns a schema instance
// the registry is given as argument in order to resolve schemas that depends on other schema.
func LoadAvroSchema(rawSchema string, registry agentiface.SchemaRegistry) (agentiface.Schema, error) {
	// retrieve all schemas from the registry that are avro
	// this is not efficient, but I can't probably do more
	ids := registry.SchemaListIds()
	schemas := make(map[string]avro.Schema)

	for _, id := range ids {
		schema, ok := registry.SchemaGetByID(id)
		if ok == nil && schema.MimeType() == MimeTypeAvroSchema {
			schemas[schema.ID()] = schema.(*AvroSchema).schema
		}
	}

	avroSchema, err := avro.ParseSchemaWithRegistry(rawSchema, schemas)

	if err != nil {
		return nil, err
	}

	title := ""
	if t, ok := avroSchema.Prop("title"); ok {
		title, _ = t.(string)
	}

	return &AvroSchema{
		id:     avroSchema.GetName(),
		title:  title,
		raw:    avroSchema.String(),
		schema: avroSchema,
	}, nil
}
