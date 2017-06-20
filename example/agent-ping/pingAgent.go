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
	"math/rand"
	"strconv"
	"time"

	"github.com/crucibuild/sdk-agent-go/agentiface"
	"github.com/crucibuild/sdk-agent-go/agentimpl"
	"github.com/crucibuild/sdk-agent-go/example"
	"io/ioutil"
	"net/http"
)

var Resources http.FileSystem

type PingAgent struct {
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

func NewPingAgent() agentiface.Agent {
	var agentSpec map[string]interface{}

	manifest := MustOpenResources("/resources/manifest.json")

	err := json.Unmarshal(manifest, &agentSpec)

	if err != nil {
		// FIXME: remove panic, and return an error -> let the caller decide what to do with that
		panic(err)
	}

	agent := &PingAgent{
		agentimpl.NewAgent(agentSpec),
	}

	if err := agent.init(); err != nil {
		panic(err)
	}

	return agent
}

func (a *PingAgent) register(rawAvroSchema string) error {
	s, err := agentimpl.LoadAvroSchema(rawAvroSchema, a)
	if err != nil {
		return err
	}

	a.SchemaRegister(s)

	return nil
}

func (a *PingAgent) init() error {
	a.SetDefaultConfigOption("delay", 1000)

	// registers additional CLI options
	for _, c := range a.Cli.RootCommand().Commands() {
		if c.Use == "start" {
			c.Flags().Int32("delay", 1000, "The delay for pinging")
			break
		}
	}

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

func (a *PingAgent) onStateChange(state agentiface.State) error {
	switch state {
	case agentiface.STATE_CONNECTED:
		// register callbacks
		_, err := a.RegisterEventCallback(map[string]interface{}{
			"type": "crucibuild/agent-example-go#tested-event",
		}, a.onTestedEvent)

		if err != nil {
			return err
		}

		// register main ping function
		a.Go(func(quit <-chan struct{}) error {
			delay, err := strconv.Atoi(a.GetConfigString("delay"))

			if err != nil {
				return err
			}

			for {
				select {
				case <-quit:
					return nil
				case <-time.After(time.Duration(delay) * time.Millisecond):
					// send command to pong agent
					cmd := &example.TestCommand{Foo: &example.Header{Z: "ok"}, Value: "ping", X: rand.Int31n(1000)}

					err := a.SendCommand("agent-pong", cmd)

					if err != nil {
						a.Error(err.Error())
					} else {
						a.Info("ping")
					}
				}
			}
		})
	case agentiface.STATE_DISCONNECTED:
		println("Disconnected!!!")
	}
	return nil

}

func (a *PingAgent) onTestedEvent(ctx agentiface.EventCtx) error {
	a.Info("Receive tested-event: " + ctx.Message().(*example.TestedEvent).Value)

	return nil
}
