package agentimpl

import (
	"github.com/crucibuild/sdk-agent-go/agentiface"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewTypeFromInterface(t *testing.T) {
	// Arrange
	i := struct {
		Value string
	}{
		"bar",
	}
	expectedName := "foo"
	expectedType := reflect.TypeOf(i)

	// Act
	tpe, err := NewTypeFromInterface(expectedName, i)

	// Assert
	assert.Nil(t, err, "No error expected")
	assert.NotNil(t, tpe, "Type must be not nil")
	assert.Equal(t, expectedName, tpe.Name(), "Names must be equal")
	assert.Equal(t, expectedType, tpe.Type(), "Types must be equal")
}

func TestNewTypeFromType(t *testing.T) {
	// Arrange
	expectedName := "foo"
	expectedType := reflect.TypeOf("")

	// Act
	tpe := NewTypeFromType(expectedName, expectedType)

	// Assert
	assert.NotNil(t, tpe, "Type must be not nil")
	assert.Equal(t, expectedName, tpe.Name(), "Names must be equal")
	assert.Equal(t, expectedType, tpe.Type(), "Types must be equal")
}

func TestNewTypeRegistry(t *testing.T) {
	// Arrange
	var agent agentiface.Agent = nil // not used

	// Act
	registry := NewTypeRegistry(agent)

	// Assert
	assert.NotNil(t, registry, "Registry instance must be not nil")
	assert.Equal(t, 0, len(registry.TypeListNames()), "Registry must be empty")
}

func TestRegisterANewType(t *testing.T) {
	// Arrange
	var agent agentiface.Agent = nil // not used
	expectedType := NewTypeFromType("foo", reflect.TypeOf(""))
	registry := NewTypeRegistry(agent)

	// Act
	registry.TypeRegister(expectedType)

	// Assert
	assert.Equal(t, 1, len(registry.TypeListNames()), "Registry must contain one type")
	assert.Equal(t, true, registry.TypeExist("foo"), "Registry must contain 'foo' type")

	tpe, err := registry.TypeGetByName("foo")
	assert.Nil(t, err, "Type must be retrieved by name ('foo')")
	assert.Equal(t, expectedType, tpe, "Types must match")

	tpe, err = registry.TypeGetByType(reflect.TypeOf(""))
	assert.Nil(t, err, "Type must be retrieved by type ('string')")
	assert.Equal(t, expectedType, tpe, "Types must match")
}

func TestUnregisterAType(t *testing.T) {
	// Arrange
	var agent agentiface.Agent = nil // not used
	expectedType := NewTypeFromType("foo", reflect.TypeOf(""))
	registry := NewTypeRegistry(agent)
	registry.TypeRegister(expectedType)

	// Act
	err := registry.TypeUnregister("foo")

	// Assert
	assert.Nil(t, err, "Unregistering type must succeed")
	assert.Equal(t, 0, len(registry.TypeListNames()), "Registry must be empty")
}
