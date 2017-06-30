package agentimpl

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	AGENT_NAME        = "the name"
	AGENT_DESCRIPTION = "short description"
	AGENT_VERSION     = "1.0.1"
	AGENT_FOO         = "bar"
)

func TestNewManifest(t *testing.T) {
	Convey(fmt.Sprintf("Given a specification for the agent '%s'", AGENT_NAME), t, func() {
		agentSpec := make(map[string]interface{})

		agentSpec["name"] = AGENT_NAME
		agentSpec["description"] = AGENT_DESCRIPTION
		agentSpec["version"] = AGENT_VERSION
		agentSpec["foo"] = AGENT_FOO

		Convey("When we create a new manifest", func() {
			manifest := NewManifest(agentSpec)

			Convey("The manifest should not be nil", func() {
				So(manifest, ShouldNotBeNil)
			})

			Convey("The manifest should be equal to the specification", func() {
				So(manifest.Name(), ShouldEqual, AGENT_NAME)
				So(manifest.Description(), ShouldEqual, AGENT_DESCRIPTION)
				So(manifest.Version(), ShouldEqual, AGENT_VERSION)
				So(manifest.Spec(), ShouldEqual, agentSpec)
			})
		})
	})
}
