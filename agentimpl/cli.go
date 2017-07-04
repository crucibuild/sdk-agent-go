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
	"github.com/spf13/cobra"
)

// Cli implements command line argument parsing.
type Cli struct {
	agent   agentiface.Agent
	rootCmd *cobra.Command
}

// NewCli creates an new Command Line Interface for the agent.
func NewCli(a agentiface.Agent) *Cli {
	cli := &Cli{
		agent: a,
		rootCmd: &cobra.Command{
			Use:   a.Name(),
			Short: a.Description(),
			Long:  "",
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				a.Info(a.Id())

				// get the flag 'config' if set (--config)
				file, _ := cmd.Root().PersistentFlags().GetString("config")

				if file == "" {
					// default configuration file
					file = fmt.Sprintf("%s/%s", fmt.Sprintf(agentiface.ConfigPathLocal, a.Name()), agentiface.ConfigName)
				}

				err := a.LoadConfigFrom(file)

				if err != nil {
					// do not stop, just report the error (default values will be used)
					a.Warning("Failed to load configuration file. %s", err.Error())
				} else {
					a.Debug("Configuration file: %s", util.AbsPathify(file))
				}
				a.Debug("Executing command <%s>", cmd.Name())
			},
			SilenceErrors: true,
		},
	}

	return cli
}

// ParseCommandLine parse the arguments
func (cli *Cli) ParseCommandLine() error {
	err := cli.rootCmd.Execute()

	if err != nil {
		cli.agent.Error("Command failed: %s", err.Error())
	}

	return err
}

// RegisterCommand register additional commands available via the command line.
func (cli *Cli) RegisterCommand(cmd *cobra.Command) {
	cli.rootCmd.AddCommand(cmd)
}

// RootCommand return the rootCommand of the agent.
func (cli *Cli) RootCommand() *cobra.Command {
	return cli.rootCmd
}
