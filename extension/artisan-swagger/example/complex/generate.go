package complex_example

import (
	"github.com/Myriad-Dreamin/artisan/artisan-core"
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

func Generate() *artisan_core.PublishedServices {

	v1 := "v1"

	userCate := DescribeUserService("")
	objectCate := DescribeObjectService("")
	objectCate.ToFile("control/object.go")

	return artisan_core.NewService(
		userCate, objectCate).Base(v1).Final()
}
