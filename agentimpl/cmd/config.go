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
	"fmt"

	"github.com/crucibuild/sdk-agent-go/agentiface"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

// RegisterCmdConfig registers command line "config" command which enables the user to interact with the agent config.
func RegisterCmdConfig(a agentiface.Agent) {
	// Manage flags:
	a.RootCommand().PersistentFlags().String("config", "", fmt.Sprintf("config file (default is $HOME/.%s/%s)", a.Name(), agentiface.ConfigName))

	// Register commands
	a.RegisterCommand(cmdConfigInit(a))
	a.RegisterCommand(cmdConfigList(a))
	a.RegisterCommand(cmdConfigGet(a))

}

func cmdConfigInit(a agentiface.Agent) *cobra.Command {
	command := &cobra.Command{
		Use:   "config:init",
		Short: "Initialize a configuration file with default values",
		Long:  `Initialize a configuration file with default values`,
		RunE: func(cmd *cobra.Command, args []string) error {
			overwrite, err := cmd.Flags().GetBool("overwrite")

			if err != nil {
				return err
			}

			if overwrite {
				return a.CreateDefaultConfigOverwrite()
			}
			return a.CreateDefaultConfig()
		},
	}

	command.Flags().Bool("overwrite", false, "Overwrite file if existing")

	return command
}

func cmdConfigList(a agentiface.Agent) *cobra.Command {
	command := &cobra.Command{
		Use:   "config:list",
		Short: "List all configuration",
		Long:  `List all configuration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.PrintConfig(os.Stdout)
		},
	}

	return command
}

func cmdConfigGet(a agentiface.Agent) *cobra.Command {
	command := &cobra.Command{
		Use:   "config:get",
		Short: "Get a value",
		Long:  `Get a value`,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 0:
				return errors.New("No configuration key provided")
			case 1:
				key := args[0]
				println(a.GetConfigString(key))
				return nil
			default:
				return errors.New("More than one configuration key provided")
			}
		},
	}

	return command
}
