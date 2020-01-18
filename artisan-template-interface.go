package artisan

import "fmt"

type ObjectType int
type TmplCtx interface {
	CurrentObjectType() ObjectType
	SetObjectType(ObjectType)
	Clone() TmplCtx

	AppendPackage(pkgPath string)
	MergePackages(pks PackageSet)
	GetPackages() PackageSet

	AppendUUID(uuid UUID) bool

	GetService() ServiceDescription
	GetCategories() []CategoryDescription
	PushCategory(cat CategoryDescription)
	PopCategory() CategoryDescription

	SetMeta(interface{})
	GetMeta() interface{}
}

type FuncTmplFac = func(objTmpl ObjTmpl, ctx TmplCtx) []FuncTmpl
type FuncTmplMiddleware = func(funcTmpl FuncTmpl, ctx TmplCtx) FuncTmpl

// func (:receiver) :FunctionName(:XParams) :ReturnType {
//     :statements
//     return
// }

type FuncTmplReceiver XParam
type FuncTmpl interface {
	fmt.Stringer
	WantFix() bool
	AllowError() bool
	GetName() string
	GetResponsibleObject() ObjTmpl
	GetReceiver() *FuncTmplReceiver
	GetParameterList() []*XParam
	GetReturnType() Type
	GetStatements() []*Statement
}

type ObjTmplMiddleware = func(template ObjTmpl, addition TmplCtx) ObjTmpl

// TmplType = {Struct, Interface}
// type :TmplName :TmplType {
//     ...:TmplFields
//            :Name :Type :Tag :SourceXParams
// }
//
// func method(:XParams)
//      XParam:
//          :Name :Type

// type TmplName TmplType(=) TmplFields[0].Type

type TmplType int
type ObjTmpl interface {
	fmt.Stringer
	GetName() string
	GetType() TmplType
	GetFields() []ObjTmplField
	GetSources() []*XParam
}

type ObjTmplField interface {
	GetName() string
	GetType() Type
	GetTag() map[string]string
	GetSource() *XParam
	SetSource(*XParam)
}
