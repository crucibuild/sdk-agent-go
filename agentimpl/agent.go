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
	"time"

	"github.com/cespare/wait"
	"github.com/crucibuild/sdk-agent-go/agentimpl/cmd"
	"github.com/crucibuild/sdk-agent-go/util"
)

type Agent struct {
	wait.Group

	*Cli
	*Config
	*SchemaRegistry
	*TypeRegistry
	*AMQP
	*Manifest
	*Logger

	id string
}

func NewAgent(agentSpec map[string]interface{}) *Agent {
	agent := &Agent{
		id: fmt.Sprintf("%s@%s#%d", agentSpec["name"].(string), util.Host(), time.Now().UnixNano()),
	}

	agent.Manifest = NewManifest(agentSpec)
	agent.Logger = NewLogger(agent)
	agent.Cli = NewCli(agent)

	agent.Config = NewConfig(agent)
	agent.SchemaRegistry = NewSchemaRegistry(agent)
	agent.TypeRegistry = NewTypeRegistry(agent)
	agent.AMQP = NewAMQP(agent)

	// register default commands
	cmd.RegisterCmdConfig(agent)
	cmd.RegisterCmdAgent(agent)
	cmd.RegisterCmdManifest(agent)
	cmd.RegisterCmdSchema(agent)

	return agent
}

func (a *Agent) Id() string {
	return a.id
}
