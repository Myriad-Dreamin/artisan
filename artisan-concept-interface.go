package artisan

import "fmt"

type GenTreeNode interface {
	GetName() string
	GetBase() string
	GetMeta() interface{}

	GenerateObjects(ts []FuncTmplFac, c TmplCtx) (objs []ObjTmpl, funcs []FuncTmpl)
	GetCategories() []CategoryDescription
	GetTmplFunctionFactory() []FuncTmplFac
	GetPackages() PackageSet
	GenerateRouterFunc(ctx TmplCtx, interfaceStyle string) (string, string)
}

type GenTreeNodeWithMethods interface {
	GetMethods() []MethodDescription
}

type ServiceDescription interface {
	GenTreeNode

	GetFilePath() string
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

type RouterMeta struct {
	RuntimeRouterMeta string
	NeedAuth          bool
}

func (r RouterMeta) GetRuntimeRouterMeta() string {
	return r.RuntimeRouterMeta
}

func (r RouterMeta) GetNeedAuth() bool {
	return r.NeedAuth
}

type ProposingService interface {
	Base(base string) ProposingService
	Meta(m interface{}) ProposingService
	UseModel(models ...*model) ProposingService
	Name(name string) ProposingService
	ToFile(fileName string) ProposingService

	GetBase() string
	GetMeta() interface{}
	GetModels() []*model
	GetName() string
	GetFilePath() string
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
	GetMethodType() MethodType
	GetAuthMeta() string
	GetName() string
	GetRequestProtocols() []SerializeObject
	GetResponseProtocols() []SerializeObject

	CreateMethodDescription(ctx *Context) MethodDescription
}

type DebuggerObject interface {
	DefiningPosition() string
}

type SerializeObject interface {
	DebuggerObject

	GetName() string
	CreateObjectDescription(ctx *Context) ObjectDescription
}

type Parameter interface {
	CreateParameterDescription(ctx *Context) ParameterDescription
}
