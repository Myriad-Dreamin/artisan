package artisan_core

import "fmt"

type GenTreeNodePropertyGetter interface {
	GetName() string
	GetBase() string
	GetMeta() interface{}
}

type GenTreeNode interface {
	GenTreeNodePropertyGetter

	GetCategories() []CategoryDescription
	GetTmplFunctionFactory() []FuncTmplFac
	GetPackages() PackageSet

	GenerateObjects(ts []FuncTmplFac, c TmplCtx) (objs []ObjTmpl, funcs []FuncTmpl)
	GenerateRouterFunc(ctx TmplCtx, interfaceStyle string) (string, string)
}

type GenTreeNodeWithMethods interface {
	GetMethods() []MethodDescription
}

type ServiceDescription interface {
	GenTreeNode
	proposingServiceGetter

	SetFilePath(fp string) ServiceDescription

	PublishAll(packageName string, opts *PublishOptions) error
	PublishObjects(packageName string, opts *PublishOptions) error
	PublishInterface(packageName string, opts *PublishOptions) error

	GenerateInterface(ctx TmplCtx, interfaceStyle string) string
}

type CategoryDescription interface {
	GenTreeNode

	GetPath() string
	GetMethods() []MethodDescription
	GetObjects() []ObjectDescription
	IterCategories(callback func(k string, v CategoryDescription) bool) bool

	SetName(n string) CategoryDescription
}

type MethodType int

type MethodDescription interface {
	GetMethodType() MethodType
	GetAuthMeta() string
	GetName() string
	GetRequests() []ObjectDescription
	GetReplies() []ObjectDescription
	GetPackages() PackageSet
}

type Type = fmt.Stringer

type ObjectDescType struct {
	ObjectDescription
}

func (o ObjectDescType) String() string {
	return o.GetName()
}

type ObjectDescription interface {
	DebuggerObject

	GetUUID() UUID
	GenObjectTmpl() ObjTmpl
	GetType() Type
	GetName() string
	GetEmbedObject() []ObjectDescription
	GetPackages() PackageSet
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
	GetPackages() PackageSet
}

// ProposingService is the Interface of VirtualService
// Getter of Virtual Service
// Get/Set
//    Faz
//    Models
//    Name
//    FilePath

type IRouterMeta interface {
	GetRuntimeRouterMeta() string
	GetNeedAuth() bool
}

type ServiceHumanInfo struct {
}

type proposingServiceGetter interface {
	GetHumanInfo() interface{}
	GetModels() []*model
	GetFilePath() string
}

type ProposingServiceGetter interface {
	GenTreeNodePropertyGetter
	proposingServiceGetter
}

type ProposingService interface {
	ProposingServiceGetter

	HumanInfo(interface{})
	Base(base string) ProposingService
	Meta(m interface{}) ProposingService
	UseModel(models ...*model) ProposingService
	Name(name string) ProposingService
	ToFile(fileName string) ProposingService
}

type CategoryGetter interface {
	GetName() string
	GetPath() string
	GetMeta() interface{}

	GetMethods() []Method
	GetWildObjects() []SerializeObject
	ForEachSubCate(func(path string, cat Category) (shouldStop bool)) error
}

// todo middleware
type Category interface {
	WithName(name string) Category
	Path(path string) Category
	SubCate(path string, cat Category) Category
	DiveIn(path string) Category
	Meta(m interface{}) Category

	RawMethod(m ...Method) Category
	Method(m MethodType, descriptions ...interface{}) Category

	AppendObject(objs ...SerializeObject) Category
	Object(descriptions ...interface{}) Category

	CategoryGetter

	CreateCategoryDescription(ctx *Context) CategoryDescription

	HelpWrapObjectXXX(skip int, descriptions ...interface{}) Category
}

type Method interface {
	GetName() string

	GetMethodType() MethodType
	GetAuthMeta() string
	GetRequestProtocols() []SerializeObject
	GetResponseProtocols() []SerializeObject

	CreateMethodDescription(ctx *Context) MethodDescription
}

type DebuggerObject interface {
	DefiningPosition() string
}

type SerializeObject interface {
	GetName() string
	DebuggerObject

	CreateObjectDescription(ctx *Context) ObjectDescription
}

type Parameter interface {
	CreateParameterDescription(ctx *Context) ParameterDescription
}
