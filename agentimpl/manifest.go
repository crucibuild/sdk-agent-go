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

type Manifest struct {
	spec map[string]interface{}
}

func NewManifest(agentSpec map[string]interface{}) *Manifest {
	return &Manifest{
		spec: agentSpec,
	}
}

func (agentSpec *Manifest) Name() string {
	return agentSpec.spec["name"].(string)
}

func (agentSpec *Manifest) Description() string {
	return agentSpec.spec["description"].(string)
}

func (agentSpec *Manifest) Version() string {
	return agentSpec.spec["version"].(string)
}

func (agentSpec *Manifest) Spec() map[string]interface{} {
	return agentSpec.spec
}
