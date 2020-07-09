package artisan

import (
	"reflect"
)

type transferClass struct {
	uuid     UUID
	name     string
	dp       string
	base     interface{}
	baseType reflect.Type
}

func (i transferClass) GetPackages() PackageSet {
	return PackageSetAppend(nil, reflect.TypeOf(i.base).PkgPath())
}

func (i transferClass) DefiningPosition() string {
	return i.dp
}

func (i transferClass) GetName() string {
	return i.name
}

func (i transferClass) CreateObjectDescription(*Context) ObjectDescription {
	return i
}

func (i transferClass) GetTypeString() string {
	return i.name
}

func (i transferClass) GetUUID() UUID {
	return i.uuid
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
