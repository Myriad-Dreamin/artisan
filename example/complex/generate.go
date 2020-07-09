package main

import (
	"fmt"
	"github.com/Myriad-Dreamin/artisan"
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
)

type CodeRawType = int

var codeField = artisan.Param("code", new(CodeRawType))
var required = artisan.Tag("binding", "required")

func main() {
	v1 := "v1"

	userCate := DescribeUserService("")
	objectCate := DescribeObjectService("")
	objectCate.ToFile("control/object.go")

	svc := artisan.NewService(
		userCate, objectCate).Base(v1).Final()

	userSvc := svc.GetService(userCate)
	delete(svc.GetServices(), userCate)

	sugar.HandlerError0(svc.PublishInterface(userSvc.SetFilePath("control/user-interface.go")))
	sugar.HandlerError0(svc.PublishObjects(userSvc.SetFilePath("control/user-objects.go")))

	err := svc.Publish()
	if err != nil {
		fmt.Println(err)
	}
}
