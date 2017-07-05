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

// Manifest represents the core properties defining an Agent.
type Manifest struct {
	spec map[string]interface{}
}

// NewManifest creates a new instance of Manifest from a map.
func NewManifest(agentSpec map[string]interface{}) *Manifest {
	return &Manifest{
		spec: agentSpec,
	}
}

// Name returns the agent name.
func (agentSpec *Manifest) Name() string {
	return agentSpec.spec["name"].(string)
}

// Description returns the agent description.
func (agentSpec *Manifest) Description() string {
	return agentSpec.spec["description"].(string)
}

// Version returns the agent version.
func (agentSpec *Manifest) Version() string {
	return agentSpec.spec["version"].(string)
}

// Spec return the complete agent specification.
func (agentSpec *Manifest) Spec() map[string]interface{} {
	return agentSpec.spec
}
