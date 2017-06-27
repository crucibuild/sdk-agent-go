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

import (
	"io"

	"github.com/spf13/pflag"
)

// Paths to look for the config file
const (
	// name of the configuration file
	CONFIG_FILENAME = "config.properties"
	// path to look for the config file in a local (user) place
	// %s stands for agent name
	CONFIG_PATH_LOCAL = "$HOME/.%s"
)

type Config interface {
	LoadConfig() error
	LoadConfigFrom(path string) error
	CreateDefaultConfig() error
	CreateDefaultConfigOverwrite() error
	PrintConfig(w io.Writer) error
	SetDefaultConfigOption(key string, value interface{})
	BindConfigPFlag(key string, flag *pflag.Flag) error
	GetConfigString(key string) string
}
