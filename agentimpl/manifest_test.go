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
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	AgentName        = "the name"
	AgentDescription = "short description"
	AgentVersion     = "1.0.1"
	AgentFoo         = "bar"
)

func TestNewManifest(t *testing.T) {
	Convey(fmt.Sprintf("Given a specification for the agent '%s'", AgentName), t, func() {
		agentSpec := make(map[string]interface{})

		agentSpec["name"] = AgentName
		agentSpec["description"] = AgentDescription
		agentSpec["version"] = AgentVersion
		agentSpec["foo"] = AgentFoo

		Convey("When we create a new manifest", func() {
			manifest := NewManifest(agentSpec)

			Convey("The manifest should not be nil", func() {
				So(manifest, ShouldNotBeNil)
			})

			Convey("The manifest should be equal to the specification", func() {
				So(manifest.Name(), ShouldEqual, AgentName)
				So(manifest.Description(), ShouldEqual, AgentDescription)
				So(manifest.Version(), ShouldEqual, AgentVersion)
				So(manifest.Spec(), ShouldEqual, agentSpec)
			})
		})
	})
}
