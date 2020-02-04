package artisan

import "fmt"
import "testing"

func TestInnerFromBigCamelToSnake(t *testing.T) {
	fmt.Println(fromBigCamelToSnake("test"))
	fmt.Println(fromBigCamelToSnake("Test"))
	fmt.Println(fromBigCamelToSnake("TestID"))
	fmt.Println(fromBigCamelToSnake("ABC"))
	fmt.Println(fromBigCamelToSnake("QwQ"))
}