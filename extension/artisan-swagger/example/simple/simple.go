package main

import (
	"fmt"
	"github.com/Myriad-Dreamin/artisan/artisan-core"
)

func main() {
	helloWorldObject := artisan_core.Object("HelloWorldObject",
		artisan_core.Param("hello_world", artisan_core.String, artisan_core.Tag("k", "v")))

	err := artisan_core.NewService().
		AppendObject(helloWorldObject, helloWorldObject).
		SetPackageName("main").WildToFile("object.go").Publish()
	if err != nil {
		fmt.Println(err)
	}
}
