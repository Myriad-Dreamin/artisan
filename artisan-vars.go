package artisan

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	VarContextRouteStructName   = "route-struct-name"
	VarContextServiceStructName = "service-struct-name"
)

var (
	ErrNotStruct       = errors.New("not struct")
	ErrConflictPath    = errors.New("conflict path")
	ErrNotMultipleOf2  = errors.New("not multiple of 2")
	ErrStopped         = errors.New("stopped")
	ErrMissingFilePath = errors.New("missing file path")
)

type ErrorObjectHasNoName struct {
	rawSvc ProposingService
	obj    SerializeObject
}

func (e ErrorObjectHasNoName) Error() string {
	return fmt.Sprintf(
		"newSerializeObject name must be specified: in service(%s %s) newSerializeObject defining path %s",
		getServiceName(e.rawSvc.GetName()), reflect.TypeOf(e.rawSvc), e.obj.DefiningPosition())
}

//noinspection GoUnusedFunction
func printCategories(descriptions []CategoryDescription) (res string) {
	for _, cat := range descriptions {
		if len(res) != 0 {
			res += ","
		}
		res += fmt.Sprintf("%s %s", cat.GetName(), cat.GetPath())
	}
	return
}

func errObjectHasNoName(obj SerializeObject, ctx *Context) error {
	return &ErrorObjectHasNoName{
		rawSvc: ctx.rawSvc,
		obj:    obj,
	}
}

//noinspection GoUnusedConst
const (
	POST MethodType = iota
	GET
	PATCH
	HEAD
	PUT
	DELETE
	OPTION
	CONNECT
	TRACE
)

var MethodTypeMapping = map[MethodType]string{
	POST:    "POST",
	GET:     "GET",
	PATCH:   "PATCH",
	HEAD:    "HEAD",
	PUT:     "PUT",
	DELETE:  "DELETE",
	OPTION:  "OPTION",
	CONNECT: "CONNECT",
	TRACE:   "TRACE",
}

const (
	ObjectTypeRequest ObjectType = iota
	ObjectTypeReply
	ObjectTypeObject
)

//noinspection GoUnusedConst
const (
	TmplTypeStruct TmplType = iota
	TmplTypeInterface
	TmplTypeEq
)
