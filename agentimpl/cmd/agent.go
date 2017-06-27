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

package cmd

import (
	"github.com/crucibuild/sdk-agent-go/agentiface"
	"github.com/spf13/cobra"
)

func RegisterCmdAgent(a agentiface.Agent) {
	// Manage flags:

	// Register commands
	a.RegisterCommand(cmdAgentStart(a))
}

func cmdAgentStart(a agentiface.Agent) *cobra.Command {
	command := &cobra.Command{
		Use:   "agent:start",
		Short: "Start the agent",
		Long:  `Start the agent.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			err := a.Connect()

			if err != nil {
				return err
			}

			return nil
		},
	}

	command.Flags().String("endpoint", "", "endpoint")
	a.BindConfigPFlag("endpoint", command.Flags().Lookup("endpoint"))

	return command
}
