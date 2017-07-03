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
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	AGENT_NAME        = "the name"
	AGENT_DESCRIPTION = "short description"
	AGENT_VERSION     = "1.0.1"
	AGENT_FOO         = "bar"
)

func TestNewManifest(t *testing.T) {
	// Arrange
	agentSpec := make(map[string]interface{})

	agentSpec["name"] = AGENT_NAME
	agentSpec["description"] = AGENT_DESCRIPTION
	agentSpec["version"] = AGENT_VERSION
	agentSpec["foo"] = AGENT_FOO

	// Act
	manifest := NewManifest(agentSpec)

	// Assert
	assert.Equal(t, AGENT_NAME, manifest.Name(), "Name mismatch")
	assert.Equal(t, AGENT_DESCRIPTION, manifest.Description(), "Description mismatch")
	assert.Equal(t, AGENT_VERSION, manifest.Version(), "Version mismatch")
	assert.ObjectsAreEqual(agentSpec, manifest.Spec())
}
