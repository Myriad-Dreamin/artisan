package artisan

import "fmt"

type GenTreeNode interface {
	GenerateObjects(ts []FuncTmplFac, c TmplCtx) (objs []ObjTmpl, funcs []FuncTmpl)
	GetCategories() []CategoryDescription
	GetTmplFunctionFactory() []FuncTmplFac
	GetPackages() PackageSet
}

type ServiceDescription interface {
	GenTreeNode

	GetName() string
	GetBase() string
	GetFilePath() string
}

type CategoryDescription interface {
	GenTreeNode

	GetName() string
	GetPath() string
	GetMethods() []MethodDescription
	GetObjects() []ObjectDescription
}

type MethodDescription interface {
	GetMethodType() MethodType
	GetName() string
	GetRequests() []ObjectDescription
	GetReplies() []ObjectDescription
}

type Type = fmt.Stringer
type ObjectDescription interface {
	GetUUID() UUID
	GenObjectTmpl() ObjTmpl
	GetType() Type
	GetEmbedObject() []ObjectDescription
}

type TagI = map[string]string
type Field = fmt.Stringer
type ParameterDescription interface {
	Hash() []byte
	GetSource() *source
	GetDTOName() string
	GetType() Type
	GetField() Field
	GetTag() TagI
	GetEmbedObjects() []ObjectDescription
}

// ProposingService is the Interface of VirtualService
// Getter of Virtual Service
// Get/Set
//    Faz
//    Models
//    Name
//    FilePath
type ProposingService interface {
	Base(base string) ProposingService
	UseModel(models ...*model) ProposingService
	Name(name string) ProposingService
	ToFile(fileName string) ProposingService

	GetBase() string
	GetModels() []*model
	GetName() string
	GetFilePath() string
}

type MethodType int

// todo middleware
type Category interface {
	Path(path string) Category
	SubCate(path string, cat Category) Category
	DiveIn(path string) Category

	RawMethod(m ...Method) Category
	Method(m MethodType, descriptions ...interface{}) Category

	Object(descriptions ...interface{}) Category
	AppendObject(objs ...SerializeObject) Category

	GetPath() string

	CreateCategoryDescription(ctx *Context) CategoryDescription

	HelpWrapObjectXXX(skip int, descriptions ...interface{}) Category
}

type Method interface {
	GetName() string

	CreateMethodDescription(ctx *Context) *methodDescription
}

type SerializeObject interface {
	DebuggerObject
	CreateObjectDescription(ctx *Context) ObjectDescription
}

type Parameter interface {
	CreateParameterDescription(ctx *Context) ParameterDescription
}

type DebuggerObject interface {
	DefiningPosition() string
}
