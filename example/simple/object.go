
package main

import (
    "github.com/Myriad-Dreamin/minimum-lib/controller"

)

var _ controller.MContext


type HelloWorldObject struct {
    HelloWorld string `form:"hello_world" k:"v" json:"hello_world"`
}
func PSerializeHelloWorldObject(_helloWorld string) *HelloWorldObject {

    return &HelloWorldObject{
        HelloWorld: _helloWorld,
    }
}
func SerializeHelloWorldObject(_helloWorld string) HelloWorldObject {

    return HelloWorldObject{
        HelloWorld: _helloWorld,
    }
}
func _packSerializeHelloWorldObject(_helloWorld string) HelloWorldObject {

    return HelloWorldObject{
        HelloWorld: _helloWorld,
    }
}
func PackSerializeHelloWorldObject(_helloWorld []string) (pack []HelloWorldObject) {
	for i := range _helloWorld {
		pack = append(pack, _packSerializeHelloWorldObject(_helloWorld[i]))
	}
	return
}
