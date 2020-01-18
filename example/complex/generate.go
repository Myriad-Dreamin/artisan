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

	userCate := DescribeUserService(v1)
	userCate.ToFile("control/user.go")
	objectCate := DescribeObjectService(v1)
	objectCate.ToFile("control/object.go")
	err := artisan.NewService(
		userCate, objectCate).Publish()
	if err != nil {
		fmt.Println(err)
	}
}
