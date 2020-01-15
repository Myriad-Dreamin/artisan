# artisan

中文文档：README-zh.md
英文文档：README.md

artisan用于生成一些琐碎的框架代码，目前版本（v0.8.0）暂时仅支持一部分功能。即将支持的功能参见todo-list。

### 开始

```go
package simple

import (
	"fmt"
	"github.com/Myriad-Dreamin/artisan"
)

func main() {
	err := artisan.NewService().
		Object("HelloWorldObject",
			artisan.Param("hello_world", artisan.String, artisan.Tag("k", "v")),
		).
		SetPackageName("main").WildPath("object.go").Publish()
	if err != nil {
		fmt.Println(err)
	}
}
```

在主函数中，由`artisan.NewService()`开始的函数链创建了一个新的`Service`模板。
Service允许直接创建Object，这些Object不依赖于任何其他概念，所以称为`wildObject`。

`artisan.Param("hello_world", artisan.String, artisan.Tag("k", "v"))`指定了一个序列化后结果为`hello_world`的域，默认要求序列名遵守蛇形命名规则，此时字段名为对应大驼峰命名规则。

`SetPackageName("main").WildPath("object.go").Publish()`明确了目标文件。

运行`go run .`得到`object.go`文件：

```go
package main

import (
    "github.com/Myriad-Dreamin/minimum-lib/controller"

)

var _ controller.MContext


type HelloWorldObject struct {
    HelloWorld string `form:"hello_world" k:"v" json:"hello_world"`
}
func PSerializeHelloWorldObject(_helloWorld string) *HelloWorldObject {

    return &HelloWorldObject{
        HelloWorld: _helloWorld,
    }
}
func SerializeHelloWorldObject(_helloWorld string) HelloWorldObject {

    return HelloWorldObject{
        HelloWorld: _helloWorld,
    }
}
func _packSerializeHelloWorldObject(_helloWorld string) HelloWorldObject {

    return HelloWorldObject{
        HelloWorld: _helloWorld,
    }
}
func PackSerializeHelloWorldObject(_helloWorld []string) (pack []HelloWorldObject) {
	for i := range _helloWorld {
		pack = append(pack, _packSerializeHelloWorldObject(_helloWorld[i]))
	}
	return
}
```

### 更复杂的例子

在`artisan`中，模板的组织顺序为`Service > Service > Category > Method > Objects`

`Service`是一组接口集合。用artisan生成一个Service的常用办法是继承VirtualService，如：
```go
package types
// 请求生成ObjectService
type ObjectCategories struct {
	//继承VirtualSerivce，避免重复定义Serivce的基础方法
    artisan.VirtualService
	//List对应v1/object-list GET方法
	List    artisan.Category
	//Post对应v1/object Post方法
	Post    artisan.Category
	//Inspect对应v1/object/:oid/inspect GET方法
	Inspect artisan.Category
	//IdGroup对应v1/ojbect/:oid的GET/PUT/DELETE方法
	IdGroup artisan.Category
}
```

如果你熟悉`gin`，那么`Category`大体相当于一个`gin.RouteGroup`，且`artisan.Method(MethodType, ...)`对应为`gin.Handle(httpMethod, ...)`

`Object`依赖于一个上下文，一般来说`Method`概念下的`Object`要么是一个请求`Request`，要么是一个响应`Response/Reply`，默认`Object`名字为`CategoryName + Request/Reply`，如果没有上下文（比如`wildObject`）或产生冲突，编译器会产生一个错误。

generate.go如下：

```go
package main

import (
	"fmt"
	"github.com/Myriad-Dreamin/artisan"
)

type CodeRawType = int

var codeField = artisan.Param("code", new(CodeRawType))
var required = artisan.Tag("binding", "required")

func main() {
	v1 := "v1"

	userCate := DescribeUserService(v1)
	userCate.ToFile("control/user.go")
	objectCate := DescribeObjectService(v1)
	objectCate.ToFile("control/object.go")
	err := artisan.NewService(
		userCate, objectCate).Publish()
	if err != nil {
		fmt.Println(err)
	}
}
```

object.go如下：

```go
package main

import (
	"github.com/Myriad-Dreamin/artisan"
	"github.com/Myriad-Dreamin/artisan/example/complex/model"
)

type ObjectCategories struct {
	artisan.VirtualService
	List    artisan.Category
	Post    artisan.Category
	Inspect artisan.Category
	IdGroup artisan.Category
}

func DescribeObjectService(base string) artisan.ProposingService {
	var objectModel = new(model.Object)
	svc := &ObjectCategories{
		List: artisan.Ink().
			Path("object-list").
			Method(artisan.POST, "ListObjects",
				artisan.QT("ListObjectsRequest", model.Filter{}),
				artisan.Reply(
					codeField,
					artisan.ArrayParam(artisan.Param("objects", objectModel)),
				),
			),
		Post: artisan.Ink().
			Path("object").
			Method(artisan.POST, "PostObject",
				artisan.Request(),
				artisan.Reply(
					codeField,
					artisan.Param("object", &objectModel),
				),
			),
		Inspect: artisan.Ink().Path("object/:oid/inspect").
			Method(artisan.GET, "InspectObject",
				artisan.Reply(
					codeField,
					artisan.Param("object", &objectModel),
				),
			),
		IdGroup: artisan.Ink().
			Path("object/:oid").
			Method(artisan.GET, "GetObject",
				artisan.Reply(
					codeField,
					artisan.Param("object", &objectModel),
				)).
			Method(artisan.PUT, "PutObject",
				artisan.Request()).
			Method(artisan.DELETE, "Delete"),
	}
	svc.Name("ObjectService").Base(base) //.
	//UseModel(serial.Model(serial.Name("object"), &objectModel))
	return svc
}
```

user.go如下：

```go
package main

import (
	"github.com/Myriad-Dreamin/artisan"
	"github.com/Myriad-Dreamin/artisan/example/complex/model"
)

type UserCategories struct {
	artisan.VirtualService
	List           artisan.Category
	Login          artisan.Category
	Register       artisan.Category
	ChangePassword artisan.Category
	Inspect        artisan.Category
	IdGroup        artisan.Category
}

func DescribeUserService(base string) artisan.ProposingService {
	var userModel = new(model.User)
	var vUserModel model.User
	svc := &UserCategories{
		List: artisan.Ink().
			Path("user-list").
			Method(artisan.POST, "ListUsers",
				artisan.QT("ListUsersRequest", model.Filter{}),
				artisan.Reply(
					codeField,
					artisan.ArrayParam(artisan.Param("users", artisan.Object(
						"ListUserReply",
						artisan.SPsC(&vUserModel.NickName, &vUserModel.LastLogin),
					))),
				),
			),
		Login: artisan.Ink().
			Path("login").
			Method(artisan.POST, "Login",
				artisan.Request(
					artisan.SPsC(&userModel.ID, &userModel.NickName, &userModel.Phone),
					artisan.Param("password", artisan.String, required),
				),
				artisan.Reply(
					codeField,
					artisan.SPsC(&userModel.ID, &userModel.Phone, &userModel.NickName, &userModel.Name),
					artisan.Param("identity", artisan.Strings),
					artisan.Param("token", artisan.String),
					artisan.Param("refresh_token", artisan.String),
				),
			),
		Register: artisan.Ink().
			Path("register").
			Method(artisan.POST, "Register",
				artisan.Request(
					artisan.SPs(artisan.C(&userModel.Name, &userModel.NickName, &userModel.Phone), required),
					artisan.Param("password", artisan.String, required),
				),
				artisan.Reply(
					codeField,
					artisan.Param("id", &userModel.ID)),
			),
		ChangePassword: artisan.Ink().
			Path("user/:id/password").
			Method(artisan.PUT, "ChangePassword",
				artisan.Request(
					artisan.Param("old_password", artisan.String, required),
					artisan.Param("new_password", artisan.String, required),
				),
			),
		Inspect: artisan.Ink().Path("user/:id/inspect").
			Method(artisan.GET, "InspectUser",
				artisan.Reply(
					codeField,
					artisan.Param("user", &userModel),
				),
			),
		IdGroup: artisan.Ink().
			Path("user/:id").
			Method(artisan.GET, "GetUser",
				artisan.Reply(
					codeField,
					artisan.SPsC(&userModel.NickName, &userModel.LastLogin),
				)).
			Method(artisan.PUT, "PutUser",
				artisan.Request(
					artisan.Param("phone", &userModel.Phone),
				)).
			Method(artisan.DELETE, "Delete"),
	}
	svc.Name("UserService").Base(base).UseModel(
		artisan.Model(artisan.Name("user"), &userModel),
		artisan.Model(artisan.Name("vUser"), &vUserModel))
	return svc
}
```

生成的文件参考`https://github.com/Myriad-Dreamin/artisan/tree/8deac2c1128ae488a5e96e284766413421132110/example/complex/control`

# artisan常量

```go
// VARIABLES

var (
	ErrNotStruct      = errors.New("not struct")
	ErrConflictPath   = errors.New("conflict path")
	ErrNotMultipleOf2 = errors.New("not multiple of 2")
)
var (
	Interface = &_interface
	Bool      = new(bool)
	String    = new(string)
	Strings   = new([]string)
	Byte      = new(byte)
	Bytes     = new([]byte)
	Rune      = new(rune)
	Runes     = new([]rune)
	Time      = new(time.Time)
	PTime     = new(*time.Time)

	Int   = new(int)
	Int8  = new(int8)
	Int16 = new(int16)
	Int32 = new(int32)
	Int64 = new(int64)

	Uint   = new(int)
	Uint8  = new(int8)
	Uint16 = new(int16)
	Uint32 = new(int32)
	Uint64 = new(int64)
)
var DefaultFunctionTmplFactories = []FuncTmplFac{
	pMethod, vMethod, packMethod,
}
```

# artisan函数

```go
func C(pack ...interface{}) []interface{}
func GenerateObjects(
	objs []ObjTmpl, funcs []FuncTmpl)
func SetDefaultFunctionTmplFactories(funcs []FuncTmplFac)
func SetDefaultMetaFactory(_metaFac func() interface{})
func SetDefaultNewTmplContext(_newTmplContext func(svc ServiceDescription) TmplCtx)
func XParamAsVar(param *XParam) string
func Inherit(name string, bases ...interface{}) *inheritClass
func Model(descriptions ...interface{}) *model
func Tag(key, value string) *tag
func T(name string, base interface{}) *transferClass
func Transfer(name string, base interface{}) *transferClass
func NewBaseFuncTmpl(wantFix bool, rObject ObjTmpl) BaseFuncTmplImpl
func Ink(_ ...interface{}) Category
func NewFuncTmpl(wantFix bool, rObject ObjTmpl) *FuncTmplImpl
func ArrayParam(param Parameter) Parameter
func P(name string, descriptions ...interface{}) Parameter
func Param(name string, descriptions ...interface{}) Parameter
func Ps(mNameTypes []interface{}, descriptions ...interface{}) (p []Parameter)
func SP(descriptions ...interface{}) Parameter
func SPs(mTypes []interface{}, descriptions ...interface{}) (p []Parameter)
func SPsC(mTypes ...interface{}) []Parameter
func SnakeParam(descriptions ...interface{}) Parameter
func NewService(rawSvc ...ProposingService) *PublishingServices
func Reply(descriptions ...interface{}) ReplyObject
func Q(desc ...interface{}) RequestObject
func QT(name string, base interface{}) RequestObject
func Request(descriptions ...interface{}) RequestObject
func Object(descriptions ...interface{}) SerializeObject
```

# artisan概念接口

```go
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

```

# artisan模板接口

```go
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

```

