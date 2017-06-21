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

package cmd

import (
	"fmt"

	"encoding/json"
	"github.com/crucibuild/sdk-agent-go/agentiface"
	"github.com/spf13/cobra"
)

func RegisterCmdManifest(a agentiface.Agent) {
	// Manage flags:

	// Register commands
	a.RegisterCommand(cmdManifestVersion(a))
	a.RegisterCommand(cmdManifestName(a))
	a.RegisterCommand(cmdManifestShow(a))
}

func cmdManifestVersion(a agentiface.Agent) *cobra.Command {
	command := &cobra.Command{
		Use:   "manifest:version",
		Short: "Provide the version of the agent",
		Long:  `Provide the version of the agent.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("%s", a.Version()))
		},
	}

	return command
}

func cmdManifestName(a agentiface.Agent) *cobra.Command {
	command := &cobra.Command{
		Use:   "manifest:name",
		Short: "Provide the name of the agent",
		Long:  `Provide the name of the agent.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("%s", a.Name()))
		},
	}

	return command
}

func cmdManifestShow(a agentiface.Agent) *cobra.Command {
	command := &cobra.Command{
		Use:   "manifest:show",
		Short: "Show the manifest of the agent",
		Long:  `Show the manifest of the agent.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			prettyFormat, _ := cmd.Flags().GetBool("pretty-print")

			var raw []byte
			var err error
			if prettyFormat {
				raw, err = json.MarshalIndent(a.Spec(), "", "  ")
			} else {
				raw, err = json.Marshal(a.Spec())
			}

			if err != nil {
				return err
			}

			fmt.Println(string(raw[:]))

			return nil
		},
	}

	command.Flags().Bool("pretty-print", false, "pretty print the manifest")

	return command
}
