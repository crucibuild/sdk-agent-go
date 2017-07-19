package agentimpl

import (
	"fmt"
	"github.com/crucibuild/sdk-agent-go/agentiface"
	"github.com/crucibuild/sdk-agent-go/util"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	schemaID    = "http://github.com/crucibuild/sdk-agent-go/schema/person-test"
	schemaTitle = "Person"
)

var personSchema = fmt.Sprintf(`{
			"id": "%s",
			"title": "%s",
			"type": "object",
			"properties": {
				"firstName": {
					"type": "string"
				},
				"lastName": {
					"type": "string"
				},
				 "age": {
            		"description": "Age in years",
            		"type": "integer",
            		"minimum": 0
        		}
			},
			"required": ["firstName", "lastName"]
		}`, schemaID, schemaTitle)

type person struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}

var personType agentiface.Type

func TestLoadJsonSchema(t *testing.T) {
	Convey(fmt.Sprintf("Given a JSon Schema (%s)", schemaTitle), t, func() {

		Convey("When when we call the LoadJSONSchema() function", func() {
			schema, err := LoadJSONSchema(personSchema)

			Convey("No error should occur", func() {
				So(err, ShouldBeNil)
			})

			Convey("Schema instance should not be nil", func() {
				So(schema, ShouldNotBeNil)
			})

			Convey(fmt.Sprintf("Schema id should equal %s", schemaID), func() {
				So(schema.ID(), ShouldEqual, schemaID)
			})

			Convey(fmt.Sprintf("Schema title should equal %s", schemaID), func() {
				So(schema.Title(), ShouldEqual, schemaTitle)
			})

			Convey(fmt.Sprintf("Schema mimetype should equal %s", MimeTypeJSONSchema), func() {
				So(schema.Title(), ShouldEqual, schemaTitle)
			})
		})
	})
}

func TestJSONDecode(t *testing.T) {
	Convey(fmt.Sprintf(`Given:
- a schema instance (mimetype %s)
- a type (%s)
- an array of bytes representing a json instance of the schema`, MimeTypeJSONSchema, personType.Name()), t, func() {
		schema, err := LoadJSONSchema(personSchema)
		payload := []byte(`{"firstName": "john", "lastName": "doe", "age": 74}`)

		Convey("No error should occur when loading schema", func() {
			So(err, ShouldBeNil)
		})

		Convey("When when we decode the bytes", func() {
			decoded, err := schema.Decode(payload, personType)

			Convey("No error should occur", func() {
				So(err, ShouldBeNil)
			})

			Convey("Decoded instance should not be nil", func() {
				So(decoded, ShouldNotBeNil)
			})

			Convey(fmt.Sprintf("Decoded instance should have the type '%s'", personType.Name()), func() {
				So(decoded, ShouldHaveSameTypeAs, &person{})
			})

			Convey("Decoded instance should have expected values", func() {
				p := decoded.(*person)

				So(p.FirstName, ShouldEqual, "john")
				So(p.LastName, ShouldEqual, "doe")
				So(p.Age, ShouldEqual, 74)
			})
		})
	})
}

func TestJSONEncode(t *testing.T) {
	Convey(fmt.Sprintf(`Given:
- a schema instance (mimetype %s)
- a type (%s)
- an instance of that type`, MimeTypeJSONSchema, personType.Name()), t, func() {
		schema, err := LoadJSONSchema(personSchema)
		p := &person{FirstName: "john", LastName: "doe", Age: 74}

		Convey("No error should occur when loading schema", func() {
			So(err, ShouldBeNil)
		})

		Convey("When when we encode the instance", func() {
			coded, err := schema.Code(p)

			Convey("No error should occur", func() {
				So(err, ShouldBeNil)
			})

			Convey("Coded bytes should not be nil", func() {
				So(coded, ShouldNotBeNil)
			})

			Convey("Coded bytes should have expected value", func() {
				fmt.Print("-->", string(coded))
				So(string(coded), ShouldEqual, `{"firstName":"john","lastName":"doe","age":74}`)
			})
		})
	})
}

func init() {
	t, err := util.GetStructType(&person{})
	if err != nil {
		panic(err)
	}
	personType = NewTypeFromType("person", t)
}
