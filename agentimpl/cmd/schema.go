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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func RegisterCmdSchema(a agentiface.Agent) {
	// Manage flags:

	// Register commands
	a.RegisterCommand(cmdSchemaList(a))
	a.RegisterCommand(cmdSchemaGet(a))
}

func cmdSchemaList(a agentiface.Agent) *cobra.Command {
	command := &cobra.Command{
		Use:   "schema:list",
		Short: "List all registered schema names",
		Long:  `List all registered schema names`,
		Run: func(cmd *cobra.Command, args []string) {
			for _, v := range a.SchemaListIds() {
				println(v)
			}
		},
	}

	return command
}

func cmdSchemaGet(a agentiface.Agent) *cobra.Command {
	command := &cobra.Command{
		Use:   "schema:get",
		Short: "Get a schema given its name",
		Long:  `Get a schema given its name`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("No schema identifier provided")
			}

			s, err := a.SchemaGetById(args[0])

			if err != nil {
				return err
			}

			println(s.Raw())

			return nil
		},
	}

	return command
}
