package main

import (
	"fmt"
	"github.com/Myriad-Dreamin/artisan"
)

type CodeRawType = int

var codeField = artisan.Param("code", new(CodeRawType))
var required = artisan.Tag("binding", "required")

func main() {
	v1 := "v1"

	userCate := DescribeUserService("")
	userCate.ToFile("control/user.go")
	objectCate := DescribeObjectService("")
	objectCate.ToFile("control/object.go")

	svc := artisan.NewService(
		userCate, objectCate).Base(v1).Final()
	err := svc.Publish()
	if err != nil {
		fmt.Println(err)
	}

}
