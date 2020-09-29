package complex_example

import (
	"github.com/Myriad-Dreamin/artisan/artisan-core"
	"github.com/go-openapi/spec"
)

type CodeRawType = int

var codeField = artisan_core.Param("code", new(CodeRawType))
var required = artisan_core.Tag("binding", "required")

type Meta struct {
	artisan_core.RouterMeta
}

func (m *Meta) NeedAuth() *Meta {
	return &Meta{
		RouterMeta: artisan_core.RouterMeta{
			RuntimeRouterMeta: m.RuntimeRouterMeta,
			NeedAuth:          true,
		},
	}
}

var UserCate = DescribeUserService("")
var ObjectCate = DescribeObjectService("").ToFile("control/object.go")
var V1 = "/v1"

func Generate() *artisan_core.PublishedServices {
	return artisan_core.NewService(
		UserCate, ObjectCate).HumanInfo(spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Info: &spec.Info{
				InfoProps: spec.InfoProps{
					Title:       "Simple Backend",
					Description: "Complex Example of Artisan",
					Version:     "v0.1.0",
				},
			},
			//ID:                  "",
			//Consumes:            nil,
			//Produces:            nil,
			//Host:                "",
			//BasePath:            "",
			//Paths:               nil,
			//Definitions:         nil,
			//Parameters:          nil,
			//Responses:           nil,
			//SecurityDefinitions: nil,
			//Security:            nil,
			//Tags:                nil,
			//ExternalDocs:        nil,
		},
	}).Base(V1).Final()
}
