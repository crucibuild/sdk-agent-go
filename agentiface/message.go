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

// State represent the current state of the agent state machine.
type State int

const (
	// ConfigDefaultEndpoint is the default endpoint for an agent.
	ConfigDefaultEndpoint = "amqp://guest:guest@localhost:5672/"

	// StateConnected represent an agent state when connected.
	StateConnected State = iota + 1

	// StateDisconnected represent an agent state when disconnected.
	StateDisconnected

	// ExchangeCommand is the name of the AMQP exchange used by the agent to send commands.
	ExchangeCommand = "crucibuild.command"

	// ExchangeEvent is the name of the AMQP exchange used by the agent to send events.
	ExchangeEvent = "crucibuild.event"

	// MimeTypeAvro is the mime type used when sending AVRO schemas.
	MimeTypeAvro = "application/vnd.apache.avro+binary"

	// AmqpHeaderSendTo is the AMQP header SendTo used to force destination of a message.
	AmqpHeaderSendTo = "SendTo"
)

// MessageName is the type representing a command name.
type MessageName string

// EventFilter is a type representing a filter on event messages.
type EventFilter map[string]interface{}

// StateCallback is a type of callback occuring on state changes.
type StateCallback func(state State) error

// Ctx denotes a context when receiving a command or an event.
// From this instance can be retrieved:
// - the message (command or event)
// - the schema
// - the properties attached to the message
type Ctx interface {
	// Messaging returns the instance of Messaging.
	Messaging() Messaging

	// Command returns the concrete instance of the deserialized message.
	Message() interface{}

	// The (Avro) schema of the message serialized.
	Schema() Schema

	// Properties returns the properties attached to the message.
	Properties() map[string]string
}

// EventCtx is a specialized Ctx which can trigger a command after receiving an event.
type EventCtx interface {
	Ctx

	// SendCommand sends a command as a consequence of this event (correlationId is set)
	SendCommand(to string, command interface{}) error
}

// CommandCtx is a specialized Ctx which can trigger a command or an event after receiving a command.
type CommandCtx interface {
	Ctx

	// SendEvent sends an event as a consequence of this message (correlationId is set)
	SendEvent(command interface{}) error

	// SendCommand sends a new command as a consequence of this command (correlationId is set)
	SendCommand(to string, command interface{}) error
}

// CommandCallback is a type of callback occuring on command reception.
type CommandCallback func(ctx CommandCtx) error

// EventCallback is a type of callback occuring on event reception.
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
