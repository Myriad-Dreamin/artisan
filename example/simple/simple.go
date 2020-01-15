package main

import (
	"fmt"
	"github.com/Myriad-Dreamin/artisan"
)

func main() {
	err := artisan.NewService().
		Object("HelloWorldObject",
			artisan.Param("hello_world", artisan.String, artisan.Tag("k", "v")),
		).
		SetPackageName("main").WildPath("object.go").Publish()
	if err != nil {
		fmt.Println(err)
	}
}
