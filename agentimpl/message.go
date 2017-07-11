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
	"github.com/crucibuild/sdk-agent-go/util"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"strings"
	"time"
)

// Ctx denotes a context when receiving a command or an event.
// From this instance can be retrieved:
// - the message (command or event)
// - the schema
// - the properties attached to the message
type Ctx struct {
	amqp   *AMQP
	data   amqp.Delivery
	schema agentiface.Schema
	msg    interface{}
}

// Messaging returns the instance of Messaging.
func (ctx *Ctx) Messaging() agentiface.Messaging {
	return ctx.amqp
}

// Message returns the concrete instance of the deserialized message.
func (ctx *Ctx) Message() interface{} {
	return ctx.msg
}

// Properties returns the properties attached to the message.
func (ctx *Ctx) Properties() map[string]string {
	p := make(map[string]string)

	p["ContentType"] = ctx.data.ContentType
	p["ContentEncoding"] = ctx.data.ContentEncoding
	//p["DeliveryMode"] = ctx.data.DeliveryMode
	//p["Priority"] = ctx.data.Priority
	p["CorrelationId"] = ctx.data.CorrelationId
	p["ReplyTo"] = ctx.data.ReplyTo
	p["Expiration"] = ctx.data.Expiration
	p["MessageId"] = ctx.data.MessageId
	//p["Timestamp"] = ctx.data.Timestamp
	p["Type"] = ctx.data.Type
	p["UserId"] = ctx.data.UserId
	p["AppId"] = ctx.data.AppId

	return p
}

// Schema The (Avro) schema of the message serialized.
func (ctx *Ctx) Schema() agentiface.Schema {
	return ctx.schema
}

// SendCommand sends a command as a consequence of this event (correlationId is set)
func (ctx *Ctx) SendCommand(to string, command interface{}) error {
	publishing, err := ctx.amqp.preparePublishing(command)

	if err != nil {
		return err
	}

	if to == "" {
		to = ctx.data.ReplyTo
	}

	publishing.Headers[agentiface.AmqpHeaderSendTo] = to
	publishing.CorrelationId = ctx.data.MessageId

	return ctx.amqp.publishCommand(publishing)
}

// SendEvent sends an event as a consequence of this message (correlationId is set)
func (ctx *Ctx) SendEvent(event interface{}) error {
	publishing, err := ctx.amqp.preparePublishing(event)

	if err != nil {
		return err
	}

	publishing.Headers[agentiface.AmqpHeaderSendTo] = ctx.data.ReplyTo
	publishing.CorrelationId = ctx.data.MessageId

	return ctx.amqp.publishEvent(publishing)
}

// AMQP is a low level handler of the AMQP connectionns, events and callbacks.
type AMQP struct {
	agent *Agent

	// amqp specifics
	connection  *amqp.Connection
	channel     *amqp.Channel
	cmdQueues   [3]amqp.Queue
	cmdChannels [3]<-chan amqp.Delivery

	// callbacks on state changes
	callbacksState map[string]agentiface.StateCallback

	// callbacks for commands
	// - key is the typename of the message (name of the schema)
	// - value is the pointer to the callback function
	callbacksCmd map[agentiface.MessageName]agentiface.CommandCallback

	// callbacks for events
	// - key is the name of the queue
	// - value is the pointer to the callback function
	callbacksEvt map[string]agentiface.EventCallback

	// aggregation channel
	aggregationChannel chan func() error
}

// NewAMQP creates a new instance of AMQP
func NewAMQP(a *Agent) *AMQP {
	a.SetDefaultConfigOption("endpoint", agentiface.ConfigDefaultEndpoint)

	return &AMQP{
		agent:              a,
		callbacksState:     make(map[string]agentiface.StateCallback),
		callbacksCmd:       make(map[agentiface.MessageName]agentiface.CommandCallback),
		callbacksEvt:       make(map[string]agentiface.EventCallback),
		aggregationChannel: nil, /* opened when connecting */
	}
}

// Connect actually connects to the AMQP broker and initializes all the exchanges and queues.
func (a *AMQP) Connect() (err error) {
	oldState := a.State()

	defer func() {
		if oldState != a.State() {
			a.notifyState()
		}
	}()

	endpoint := a.agent.GetConfigString("endpoint")

	if a.State() == agentiface.StateConnected {
		err = fmt.Errorf("Already connected to: %s", endpoint)
		return
	}

	a.agent.Info("Connecting to %s", endpoint)

	a.connection, err = amqp.Dial(endpoint)
	if err != nil {
		a.Disconnect() // nolint: errcheck, silently disconnect and do not report any errors
		return
	}

	a.channel, err = a.connection.Channel()

	if err != nil {
		a.Disconnect() // nolint: errcheck, silently disconnect and do not report any errors
		return
	}

	err = a.declareExchanges()
	if err != nil {
		a.Disconnect() // nolint: errcheck, silently disconnect and do not report any errors
		return
	}

	err = a.declareQueues()
	if err != nil {
		a.Disconnect() // nolint: errcheck, silently disconnect and do not report any errors
		return
	}

	a.agent.Go(a.readMessages)

	return
}

func (a *AMQP) readMessages(quit <-chan struct{}) error {
	endpoint := a.agent.GetConfigString("endpoint")
	a.agent.Info("Connected to: %s as %s", endpoint, a.agent.ID())

	// create an aggregation channel which gathers all incoming messages
	a.aggregationChannel = make(chan func() error)

	defer func() {
		close(a.aggregationChannel)
	}()

	a.agent.Go(a.processMessages)

	for {
		select {
		case d, ok := <-a.cmdChannels[0]:
			if !ok {
				return nil
			}

			a.aggregationChannel <- func() error {
				return a.handleCommand(d)
			}
			continue
		case d, ok := <-a.cmdChannels[1]:
			if !ok {
				return nil
			}

			a.aggregationChannel <- func() error {
				return a.handleCommand(d)
			}
			continue
		case d, ok := <-a.cmdChannels[2]:
			if !ok {
				return nil
			}
			a.aggregationChannel <- func() error {
				return a.handleCommand(d)
			}
			continue
		case <-quit:
			return nil
		}
	}
}

func (a *AMQP) processMessages(quit <-chan struct{}) error {
	defer func() {
		err := a.Disconnect()
		if err != nil {
			a.agent.Warning("%s", err.Error())
		}
	}()
	for {
		select {
		case f, ok := <-a.aggregationChannel:
			if !ok {
				return nil
			}
			// call process function
			err := f()

			if err != nil {
				a.agent.Error("%s", err.Error())
			}

		case <-quit:
			return nil
		}
	}
}

// Disconnect disconnect from the broker.
func (a *AMQP) Disconnect() error {
	oldState := a.State()

	endpoint := a.agent.GetConfigString("endpoint")

	defer func() {
		a.connection = nil
		a.channel = nil
		a.aggregationChannel = nil

		if oldState != a.State() {
			a.notifyState()
		}
	}()

	if a.State() == agentiface.StateDisconnected {
		return fmt.Errorf("Not connected")
	}

	err := a.connection.Close()

	a.agent.Info("Disconnected from: %s", endpoint)

	return err
}

// State returns the state of the connection.
func (a *AMQP) State() agentiface.State {
	if a.connection == nil {
		return agentiface.StateDisconnected
	}

	return agentiface.StateConnected
}

func (a *AMQP) notifyState() {
	s := a.State()
	for _, f := range a.callbacksState {
		f(s) // nolint: errcheck, error returned by callback function is not currently used in the sdk
	}
}

func (a *AMQP) declareExchanges() (err error) {
	err = a.channel.ExchangeDeclare(
		agentiface.ExchangeCommand, // name
		"headers",                  // type
		true,                       // durable
		false,                      // auto-deleted
		false,                      // internal
		false,                      // no-wait
		nil,                        // arguments
	) // args Table

	if err != nil {
		return err
	}

	err = a.channel.ExchangeDeclare(
		agentiface.ExchangeEvent, // name
		"headers",                // type
		true,                     // durable
		false,                    // auto-deleted
		false,                    // internal
		false,                    // no-wait
		nil,                      // arguments
	) // args Table

	return err
}

// declareQueues perform the declaration and binding of common queues for the agent.
// For commands (and requests) the following queues are created (if needed)
// - crucibuild/agent-git@localhost#352
// - crucibuild/agent-git@192.168.4.2
// - crucibuild/agent-git*/
func (a *AMQP) declareQueues() (err error) {
	// FIXME: find a better way to code this - ugly code
	// ## Declare queues for commands
	a.cmdQueues[0], err = a.channel.QueueDeclare(
		a.agent.ID(),
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return err
	}

	err = a.channel.QueueBind(
		a.cmdQueues[0].Name, // queue name
		"",                  // routing key
		agentiface.ExchangeCommand, // exchange
		false,
		amqp.Table{
			agentiface.AmqpHeaderSendTo: a.cmdQueues[0].Name,
		},
	)

	if err != nil {
		return err
	}

	err = a.channel.QueueBind(
		a.cmdQueues[0].Name, // queue name
		"",                  // routing key
		agentiface.ExchangeCommand, // exchange
		false,
		amqp.Table{
			agentiface.AmqpHeaderSendTo: "*",
		},
	)

	if err != nil {
		return err
	}

	a.cmdQueues[1], err = a.channel.QueueDeclare(
		fmt.Sprintf("%s@%s", a.agent.Manifest().Name(), util.Host()),
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return err
	}

	err = a.channel.QueueBind(
		a.cmdQueues[1].Name, // queue name
		"",                  // routing key
		agentiface.ExchangeCommand, // exchange
		false,
		amqp.Table{
			agentiface.AmqpHeaderSendTo: a.cmdQueues[1].Name,
		},
	)

	if err != nil {
		return err
	}

	a.cmdQueues[2], err = a.channel.QueueDeclare(
		a.agent.Manifest().Name(),
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return err
	}

	err = a.channel.QueueBind(
		a.cmdQueues[2].Name, // queue name
		"",                  // routing key
		agentiface.ExchangeCommand, // exchange
		false,
		amqp.Table{
			agentiface.AmqpHeaderSendTo: a.cmdQueues[2].Name,
		},
	)

	if err != nil {
		return err
	}

	// - create channels on that queues
	for i := 0; i < 3; i++ {
		a.cmdChannels[i], err = a.channel.Consume(
			a.cmdQueues[i].Name, // queue
			"",                  // consumer
			true,                // auto-ack
			false,               // exclusive
			false,               // no-local
			false,               // no-wait
			nil,                 // args
		)

		if err != nil {
			return err
		}
	}

	// ## declare queues for events
	// done on demand

	return nil
}

func (a *AMQP) getSchema(d amqp.Delivery) (agentiface.Schema, error) {
	// check content type
	if d.ContentType != agentiface.MimeTypeAvro {
		return nil, fmt.Errorf("Not Acceptable: Content-type: %s", d.ContentType)
	}

	// check message type
	messageType := strings.TrimSpace(d.Type)

	if messageType == "" {
		return nil, fmt.Errorf("Not Acceptable: No Message-type provided")
	}

	s, err := a.agent.SchemaGetByID(messageType)
	if err != nil {
		println(err.Error())
		return nil, fmt.Errorf("Not Acceptable: Message-type '%s' is unknown", messageType)
	}

	return s, nil
}

// decode decodes the given delivery and returns the Avro schema, a pointer to the decoded record and eventually
// the error if something went wrong
func (a *AMQP) decode(d amqp.Delivery) (agentiface.Schema, interface{}, error) {
	messageType := strings.TrimSpace(d.Type)

	s, err := a.getSchema(d)
	if err != nil {
		return nil, nil, err
	}

	t, err := a.agent.TypeGetByName(messageType)
	if err != nil {
		return nil, nil, fmt.Errorf("Not Acceptable: Message-type '%s' is unknown", messageType)
	}

	decodedRecord, err := s.Decode(d.Body, t)

	if err != nil {
		return nil, nil, err
	}

	return s, decodedRecord, nil
}

func (a *AMQP) handleCommand(d amqp.Delivery) error {
	s, decodedRecord, err := a.decode(d)

	if err != nil {
		return err
	}

	// Dispatch to the suitable callback
	c, ok := a.callbacksCmd[agentiface.MessageName(s.ID())]

	if !ok {
		return fmt.Errorf("Not Acceptable: Message-type '%s' is not handled", s.ID())
	}

	// Invoke the callback
	err = c(&Ctx{
		amqp:   a,
		data:   d,
		schema: s,
		msg:    decodedRecord,
	})

	return err
}

func (a *AMQP) handleEvent(d amqp.Delivery, callback agentiface.EventCallback) error {
	s, decodedRecord, err := a.decode(d)

	if err != nil {
		return err
	}

	// Invoke the callback
	err = callback(&Ctx{
		amqp:   a,
		data:   d,
		schema: s,
		msg:    decodedRecord,
	})

	return err
}

// RegisterStateCallback registers a callback triggered by state changes.
func (a *AMQP) RegisterStateCallback(stateCallback agentiface.StateCallback) string {
	key := uuid.NewV4().String()
	a.callbacksState[key] = stateCallback

	return key
}

// RegisterCommandCallback registers a callback triggered by a command reception.
func (a *AMQP) RegisterCommandCallback(commandName agentiface.MessageName, commandCallback agentiface.CommandCallback) (string, error) {
	if a.State() != agentiface.StateConnected {
		return "", fmt.Errorf("Cannot register command callback if not connected")
	}

	a.callbacksCmd[commandName] = commandCallback

	return string(commandName), nil
}

// RegisterEventCallback registers a callback triggered by an event reception.
func (a *AMQP) RegisterEventCallback(filter agentiface.EventFilter, eventCallback agentiface.EventCallback) (string, error) {
	if a.State() != agentiface.StateConnected {
		return "", fmt.Errorf("Cannot register event callback if not connected")
	}

	// declare new queue (exclusive)
	queue, err := a.channel.QueueDeclare(
		"",    // FIXME: name - automatically generated
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return "", err
	}

	err = a.channel.QueueBind(
		queue.Name, // queue name
		"",         // routing key
		agentiface.ExchangeEvent, // exchange
		false, // no-Wait
		amqp.Table(filter),
	)

	if err != nil {
		return "", err
	}

	// listen to it:
	channel, err := a.channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return "", err
	}

	a.agent.Go(func(quit <-chan struct{}) error {
		for {
			select {
			case d, ok := <-channel:
				if !ok {
					return nil
				}
				a.aggregationChannel <- func() error {
					return a.handleEvent(d, eventCallback)
				}
			case <-quit:
				return nil
			}
		}
	})

	return queue.Name, nil
}

func (a *AMQP) preparePublishing(msg interface{}) (*amqp.Publishing, error) {
	// find message name from type
	atype, err := util.GetStructType(msg)

	if err != nil {
		return nil, err
	}

	typeInfo, err := a.agent.TypeGetByType(atype)

	if err != nil {
		return nil, fmt.Errorf("Not Acceptable: %s", err.Error())
	}

	// from message name get the schema:
	schema, err := a.agent.SchemaGetByID(typeInfo.Name())

	if err != nil {
		return nil, fmt.Errorf("Not Acceptable: Message-type '%s' is not handled", typeInfo.Name())
	}

	bytes, err := schema.Code(msg)

	if err != nil {
		return nil, err
	}

	// send command:
	return &amqp.Publishing{
		Timestamp:   time.Now(),
		ContentType: agentiface.MimeTypeAvro,
		MessageId:   uuid.NewV4().String(),
		Type:        schema.ID(),
		ReplyTo:     a.agent.ID(),

		Headers: map[string]interface{}{
			// used for headers routing
			"type": schema.ID(),
		},

		Body: bytes,
	}, nil
}

func (a *AMQP) publishCommand(publishing *amqp.Publishing) error {
	if a.State() != agentiface.StateConnected {
		return errors.New("Not connected")
	}

	return a.channel.Publish(
		agentiface.ExchangeCommand,
		"",
		false, // mandatory
		false, // immediate
		*publishing)
}

func (a *AMQP) publishEvent(publishing *amqp.Publishing) error {
	a.agent.Debug("Sending event")
	if a.State() != agentiface.StateConnected {
		return errors.New("Not connected")
	}

	return a.channel.Publish(
		agentiface.ExchangeEvent,
		"",
		false, // mandatory
		false, // immediate
		*publishing)
}

// SendCommand sends a command to a specific agent.
func (a *AMQP) SendCommand(to string, command interface{}) error {
	publishing, err := a.preparePublishing(command)

	if err != nil {
		return err
	}

	publishing.Headers[agentiface.AmqpHeaderSendTo] = to

	return a.publishCommand(publishing)
}
