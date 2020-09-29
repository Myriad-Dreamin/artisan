package artisan_core

import (
	"reflect"
)

func NewService(rawSvc ...ProposingService) *PublishingServices {
	return &PublishingServices{
		rawSvc:      rawSvc,
		packageName: "control",
		wildSvc:     newWildService(),
	}
}

func Ink(_ ...interface{}) Category {
	return newCategory()
}

func Inherit(name string, bases ...interface{}) SerializeObject {
	return &inheritClass{name: name, bases: bases,
		uuid: MustUUID(),
		dp:   getCaller(1).String()}
}

func Transfer(name string, base interface{}) SerializeObject {
	return &transferClass{
		name: name, base: base, baseType: reflect.TypeOf(base),
		uuid: MustUUID(),
		dp:   getCaller(1).String()}
}

func Reply(descriptions ...interface{}) ReplyObject {
	return ReplyObject{s: newSerializeObject(2, descriptions...)}
}

func Request(descriptions ...interface{}) RequestObject {
	return RequestObject{s: newSerializeObject(2, descriptions...)}
}

func Object(descriptions ...interface{}) SerializeObject {
	return newSerializeObject(1, descriptions...)
}

func Param(name string, descriptions ...interface{}) Parameter {
	return createNorm(name, descriptions...)
}

func SnakeParam(descriptions ...interface{}) Parameter {
	return newSnake(createNorm("_snaking", descriptions...))
}

func ArrayParam(param Parameter) Parameter {
	return arrayParam{p: param}
}

func Tag(key, value string) *tag {
	return &tag{
		key:   key,
		value: value,
	}
}

func NewBaseFuncTmpl(wantFix bool, rObject ObjTmpl) BaseFuncTmplImpl {
	return BaseFuncTmplImpl{
		Fix:     wantFix,
		RObject: rObject,
	}
}

func NewFuncTmpl(wantFix bool, rObject ObjTmpl) *FuncTmplImpl {
	return &FuncTmplImpl{
		BaseFuncTmplImpl: NewBaseFuncTmpl(wantFix, rObject),
	}
}
