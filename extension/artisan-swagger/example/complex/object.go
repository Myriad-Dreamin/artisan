package complex_example

import (
	"github.com/Myriad-Dreamin/artisan/artisan-core"
	"github.com/Myriad-Dreamin/artisan/extension/artisan-swagger/example/complex/model"
)

type ObjectCategories struct {
	artisan_core.VirtualService
	List    artisan_core.Category
	Post    artisan_core.Category
	Inspect artisan_core.Category
	IdGroup artisan_core.Category
}

func DescribeObjectService(base string) artisan_core.ProposingService {
	meta := Meta{
		RouterMeta: artisan_core.RouterMeta{RuntimeRouterMeta: "object"},
	}

	var objectModel = new(model.Object)
	svc := &ObjectCategories{
		List: artisan_core.Ink().
			Path("object-list").
			Method(artisan_core.POST, "ListObjects", artisan_core.AuthMeta("object@r"),
				artisan_core.QT("ListObjectsRequest", model.Filter{}),
				artisan_core.Reply(
					codeField,
					artisan_core.ArrayParam(artisan_core.Param("objects", objectModel)),
				),
			),
		Post: artisan_core.Ink().
			Path("object").
			Method(artisan_core.POST, "PostObject",
				artisan_core.Request(),
				artisan_core.Reply(
					codeField,
					artisan_core.Param("object", &objectModel),
				),
			),
		Inspect: artisan_core.Ink().Path("object/:oid/inspect").
			Method(artisan_core.GET, "InspectObject",
				artisan_core.Reply(
					codeField,
					artisan_core.Param("object", &objectModel),
				),
			),
		IdGroup: artisan_core.Ink().
			Path("object/:oid").
			Method(artisan_core.GET, "GetObject",
				artisan_core.Reply(
					codeField,
					artisan_core.Param("object", &objectModel),
				)).
			Method(artisan_core.PUT, "PutObject",
				artisan_core.Request()).
			Method(artisan_core.DELETE, "Delete"),
	}
	svc.Name("ObjectService").Base(base).Meta(meta) //.
	//UseModel(serial.Model(serial.Name("object"), &objectModel))
	return svc
}
