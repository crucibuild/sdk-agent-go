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

type State int

const (
	CONFIG_DEFAULT_ENDPOINT = "amqp://guest:guest@localhost:5672/"

	STATE_CONNECTED State = iota + 1
	STATE_DISCONNECTED

	EXCHANGE_COMMAND = "crucibuild.command"
	EXCHANGE_EVENT   = "crucibuild.event"

	MIMETYPE_AVRO = "application/vnd.apache.avro+binary"

	AMQP_HEADER_SEND_TO = "SendTo"
)

type MessageName string

type EventFilter map[string]interface{}

type StateCallback func(state State) error

// Ctx denotes a context when receiving a command or an event.
// From this instance:
// - the message (command or event) can be retrieved
// - the schema
// - the properties attached to the message
type Ctx interface {
	// Messaging returns the instance of Messaging
	Messaging() Messaging

	// Command returns the concrete instance of the deserialized message
	Message() interface{}

	// The (avro) schema of the message serialized
	Schema() Schema

	// Properties returns the properties attached to the message
	Properties() map[string]string
}

type EventCtx interface {
	Ctx

	// SendCommand sends a command as a consequence of this event (correlationId is set)
	SendCommand(to string, command interface{}) error
}

type CommandCtx interface {
	Ctx

	// SendEvent sends an event as a consequence of this message (correlationId is set)
	SendEvent(command interface{}) error

	// SendCommand sends a new command as a consequence of this command (correlationId is set)
	SendCommand(to string, command interface{}) error
}

type CommandCallback func(ctx CommandCtx) error

type EventCallback func(ctx EventCtx) error

// Messaging interface denotes the capability to send and receive messages and manage connection.
type Messaging interface {
	Connect() error
	Disconnect() error
	State() State

	RegisterStateCallback(stateCallback StateCallback) string

	RegisterCommandCallback(commandName MessageName, commandCallback CommandCallback) (string, error)

	RegisterEventCallback(filter EventFilter, eventCallback EventCallback) (string, error)

	SendCommand(to string, command interface{}) error
}
