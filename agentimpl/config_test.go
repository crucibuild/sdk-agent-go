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
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewConfig(t *testing.T) {
	Convey(fmt.Sprintf("Given an empty agent and the default config path."), t, func() {
		var agent agentiface.Agent // not used

		Convey(fmt.Sprintf("When when we call the NewConfig() function"), func() {
			config := NewConfig(agent)

			Convey("Configuration instance must not be nil", func() {
				So(config, ShouldNotBeNil)
			})
		})
	})
}
