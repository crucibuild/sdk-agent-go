/* Copyright (C) 2016 Christophe Camel, Jonathan Pigrée
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/json"
	"fmt"
	"github.com/crucibuild/sdk-agent-go/agentiface"
	"github.com/crucibuild/sdk-agent-go/agentimpl"
	"github.com/crucibuild/sdk-agent-go/example"
	"io/ioutil"
	"net/http"
)

var Resources http.FileSystem

type PongAgent struct {
	*agentimpl.Agent
}

func MustOpenResources(path string) []byte {
	file, err := Resources.Open(path)

	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadAll(file)

	if err != nil {
		panic(err)
	}

	return content
}

func NewPongAgent() agentiface.Agent {
	var agentSpec map[string]interface{}

	manifest := MustOpenResources("/resources/manifest.json")

	err := json.Unmarshal(manifest, &agentSpec)

	if err != nil {
		// FIXME: remove panic, and return an error -> let the caller decide what to do with that
		panic(err)
	}

	agent := &PongAgent{
		agentimpl.NewAgent(agentSpec),
	}

	if err := agent.init(); err != nil {
		panic(err)
	}

	return agent
}

func (a *PongAgent) register(rawAvroSchema string) error {
	s, err := agentimpl.LoadAvroSchema(rawAvroSchema, a)
	if err != nil {
		return err
	}

	a.SchemaRegister(s)

	return nil
}

func (a *PongAgent) init() error {
	// register schemas:
	var content []byte
	content = MustOpenResources("/schema/header.avro")
	if err := a.register(string(content[:])); err != nil {
		return err
	}

	content = MustOpenResources("/schema/test-command.avro")
	if err := a.register(string(content[:])); err != nil {
		return err
	}

	content = MustOpenResources("/schema/tested-event.avro")
	if err := a.register(string(content[:])); err != nil {
		return err
	}

	// register types
	// register types
	if _, err := a.TypeRegister(agentimpl.NewTypeFromType("crucibuild/agent-example-go#tested-event", example.TestedEventType)); err != nil {
		return err
	}
	if _, err := a.TypeRegister(agentimpl.NewTypeFromType("crucibuild/agent-example-go#test-command", example.TestCommandType)); err != nil {
		return err
	}

	// register state callback
	a.RegisterStateCallback(a.onStateChange)

	return nil
}

func (a *PongAgent) onStateChange(state agentiface.State) error {
	switch state {
	case agentiface.STATE_CONNECTED:
		if _, err := a.RegisterCommandCallback("crucibuild/agent-example-go#test-command", a.onTestCommand); err != nil {
			return err
		}
	}
	return nil

}

func (a *PongAgent) onTestCommand(ctx agentiface.CommandCtx) error {
	cmd := ctx.Message().(*example.TestCommand)

	a.Info(fmt.Sprintf("Received test-command: '%s' '%s' '%d' ", cmd.Foo.Z, cmd.Value, cmd.X))

	// reply with a tested event

	return ctx.SendEvent(&example.TestedEvent{Value: "pong"})
}
