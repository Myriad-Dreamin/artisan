package artisan

import (
	"reflect"
)

type transferClass struct {
	name     string
	base     interface{}
	baseType reflect.Type
}

func (i transferClass) CreateObjectDescription(ctx *Context) ObjectDescription {
	ctx.appendPackage(reflect.TypeOf(i.base).PkgPath())
	return i
}

func (i transferClass) GetTypeString() string {
	return i.name
}

func (i transferClass) GenObjectTmpl() ObjTmpl {
	x := &XParam{
		Name:   fromBigCamelToSnake(getReflectTypeElementType(i.baseType).Name()),
		TypeOf: i.baseType.String(),
		Source: nil,
	}
	f := &ObjectTmplFieldImpl{
		Name:   "",
		PType:  i.baseType,
		Tag:    nil,
		Source: x,
	}
	return &ObjTmplImpl{
		Name:   i.name,
		TType:  TmplTypeEq,
		Fields: []ObjTmplField{f},
		Xps:    []*XParam{x},
	}
}

func (i transferClass) GetType() Type {
	return i.baseType
}

func (i transferClass) GetEmbedObject() []ObjectDescription {
	return nil
}
