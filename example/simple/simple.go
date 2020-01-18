package main

import (
	"fmt"
	"github.com/Myriad-Dreamin/artisan"
)

func main() {
	helloWorldObject := artisan.Object("HelloWorldObject",
		artisan.Param("hello_world", artisan.String, artisan.Tag("k", "v")))

	err := artisan.NewService().
		AppendObject(helloWorldObject, helloWorldObject).
		SetPackageName("main").WildPath("object.go").Publish()
	if err != nil {
		fmt.Println(err)
	}
}
