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
