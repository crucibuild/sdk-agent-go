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
package agentimpl

import (
	"bufio"
	"fmt"
	"github.com/crucibuild/sdk-agent-go/agentiface"
	"github.com/crucibuild/sdk-agent-go/util"
	"github.com/magiconair/properties"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	agent agentiface.Agent
	viper *viper.Viper
}

func NewConfig(a agentiface.Agent) *Config {
	configExtension := filepath.Ext(agentiface.CONFIG_FILENAME)
	configBaseName := strings.TrimSuffix(agentiface.CONFIG_FILENAME, configExtension)

	// initialize viper (for configuration management)
	viper := viper.New()
	viper.SetConfigName(configBaseName)
	viper.SetConfigType(strings.TrimLeft(configExtension, "."))

	return &Config{
		agent: a,
		viper: viper,
	}
}

func (config *Config) SetDefaultConfigOption(key string, value interface{}) {
	config.viper.SetDefault(key, value)
}

func (config *Config) BindConfigPFlag(key string, flag *pflag.Flag) error {
	return config.viper.BindPFlag(key, flag)
}

func (config *Config) LoadConfig() error {
	err := config.viper.ReadInConfig()

	return err
}

func (config *Config) LoadConfigFrom(path string) error {
	// default configuration file
	config.viper.SetConfigFile(util.AbsPathify(path))
	return config.LoadConfig()
}

// Init initialize a new configuration file with default values
// returns an error if configuration file already exists
func (config *Config) CreateDefaultConfig() error {
	return config.createDefaultConfig(false)
}

// Init initialize a new configuration file with default values
// Overwrite the file if it exists
func (config *Config) CreateDefaultConfigOverwrite() error {
	return config.createDefaultConfig(true)
}

func (config *Config) createDefaultConfig(overwrite bool) error {
	cfgFile := config.viper.ConfigFileUsed()

	if util.Exists(cfgFile) {
		if overwrite {
			err := os.Remove(cfgFile)

			if err != nil {
				return err
			}
		} else {
			return errors.Errorf("Can't initialize configuration. A configuration file already exists: %s", cfgFile)
		}
	}

	p := properties.NewProperties()
	for _, v := range config.viper.AllKeys() {
		p.Set(v, config.viper.GetString(v))
	}

	err := os.MkdirAll(filepath.Dir(cfgFile), os.ModePerm)

	if err != nil {
		return err
	}

	f, err := os.Create(cfgFile)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = p.Write(f, properties.UTF8)

	if err != nil {
		return err
	}

	config.agent.Info("Configuration initialized: %s", cfgFile)

	return nil
}

func (config *Config) PrintConfig(w io.Writer) error {
	buffW := bufio.NewWriter(w)
	for _, v := range config.viper.AllKeys() {
		_, err := buffW.WriteString(fmt.Sprintf("%s=%s\n", v, config.viper.GetString(v)))

		if err != nil {
			return err
		}
	}
	err := buffW.Flush()

	return err
}

func (config *Config) GetConfigString(key string) string {
	return config.viper.GetString(key)
}
