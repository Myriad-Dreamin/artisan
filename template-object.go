package artisan

import "fmt"

type ObjTmplImpl struct {
	Name   string
	TType  TmplType
	Fields []ObjTmplField
	Xps    []*XParam
}

func (o ObjTmplImpl) String() string {
	switch o.TType {
	case TmplTypeEq:
		return fmt.Sprintf(`
type %s = %s`, o.Name, o.Fields[0].GetType().String())
	case TmplTypeStruct:
		return fmt.Sprintf(`
type %s struct {
%s
}`, o.Name, instantiateStructFields(o.Fields))
	default:
		panic("todo")
	}
}

func (o ObjTmplImpl) GetName() string {
	return o.Name
}

func (o ObjTmplImpl) GetSources() []*XParam {
	return o.Xps
}

func (o ObjTmplImpl) GetType() TmplType {
	return o.TType
}

func (o ObjTmplImpl) GetFields() []ObjTmplField {
	return o.Fields
}
