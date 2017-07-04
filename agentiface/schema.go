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

package agentiface

// Schema represent a type of message for a schema based utility like AVRO.
type Schema interface {
	// Id returns the Unique identifier used to reference the schema.
	ID() string
	// Title returns the schema title.
	Title() string
	// MimeType returns the mimetype of the schema.
	MimeType() string
	// Raw returns the schema in raw format.
	Raw() string

	// Decode unserializes a message with its schema.
	Decode(o []byte, t Type) (interface{}, error)
	// Code serializes a message with its schema.
	Code(o interface{}) ([]byte, error)
}

// SchemaRegistry permits to register multiple schemas and manage them.
type SchemaRegistry interface {
	// SchemaRegister registers a new Schema in the registry.
	SchemaRegister(schema Schema) (string, error)
	// SchemaGetByID retrieve the Schema in the registry given its id.
	SchemaGetByID(id string) (Schema, error)
	// SchemaListNames list all the Schemas contained in the registry.
	SchemaListIds() []string
	// SchemaUnregister removes the schema from the registry.
	SchemaUnregister(id string) error
	// SchemaExist check if the given schema exists in the registry.
	SchemaExist(id string) bool
}
