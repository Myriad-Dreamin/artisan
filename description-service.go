package artisan

import (
	"fmt"
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
	"os"
	"strings"
)

type serviceDescription struct {
	name          string
	base          string
	meta          interface{}
	tmplFactories []FuncTmplFac
	categories    []CategoryDescription
	filePath      string
	//packages   map[string]int
}

func (svc serviceDescription) GetName() string {
	return svc.name
}

func (svc serviceDescription) GetBase() string {
	return svc.base
}

func (svc serviceDescription) GetMeta() interface{} {
	return svc.meta
}

func (svc serviceDescription) GetTmplFunctionFactory() []FuncTmplFac {
	return svc.tmplFactories
}

func (svc serviceDescription) GetCategories() []CategoryDescription {
	return svc.categories
}

func (svc serviceDescription) GetFilePath() string {
	return svc.filePath
}

func (svc serviceDescription) GetPackages() PackageSet {
	return nil
}

func (svc *serviceDescription) SetFilePath(fp string) ServiceDescription {
	svc.filePath = fp
	return svc
}

func (svc *serviceDescription) GenerateObjects(ts []FuncTmplFac, c TmplCtx) (objs []ObjTmpl, functions []FuncTmpl) {
	return GenerateObjects(svc, ts, c)
}

func (svc *serviceDescription) PublishAll(packageName string, opts *PublishOptions) (err error) {
	if len(svc.GetFilePath()) == 0 {
		return ErrMissingFilePath
	}
	if opts == nil {
		opts = &PublishOptions{}
	}
	ctx := newTmplContext(svc)
	objectsGen := svc.GenerateObjectString(svc.GenerateObjects(DefaultFunctionTmplFactories, ctx))
	interfaceGen := svc.GenerateInterface(ctx, opts.InterfaceStyle)
	return svc.publishCtx(packageName, ctx, interfaceGen+objectsGen)
}

func (svc *serviceDescription) PublishInterface(packageName string, opts *PublishOptions) (err error) {
	if len(svc.GetFilePath()) == 0 {
		return ErrMissingFilePath
	}
	if opts == nil {
		opts = &PublishOptions{}
	}
	ctx := newTmplContext(svc)
	interfaceGen := svc.GenerateInterface(ctx, opts.InterfaceStyle)
	return svc.publishCtx(packageName, ctx, interfaceGen)
}

func (svc *serviceDescription) PublishObjects(packageName string, opts *PublishOptions) (err error) {
	if len(svc.GetFilePath()) == 0 {
		return ErrMissingFilePath
	}
	if opts == nil {
		opts = &PublishOptions{}
	}
	ctx := newTmplContext(svc)
	objectsGen := svc.GenerateObjectString(svc.GenerateObjects(DefaultFunctionTmplFactories, ctx))
	return svc.publishCtx(packageName, ctx, objectsGen)
}

func (svc *serviceDescription) publishCtx(packageName string, ctx TmplCtx, body string) (err error) {
	sugar.WithWriteFile(func(f *os.File) {
		_, err = fmt.Fprintf(f, `
package %s

import (
%s
)
%s`, packageName, depList(ctx.GetPackages()), body)
		if err != nil {
			return
		}
	}, svc.GetFilePath())
	return
}

func (svc *serviceDescription) GenerateInterface(
	ctx TmplCtx, interfaceStyle string) string {
	switch interfaceStyle {
	case InterfaceStyleMinimum:
		fallthrough
	default:
		ctx.AppendPackage("github.com/Myriad-Dreamin/minimum-lib/controller")
		if len(svc.GetName()) == 0 {
			return ""
		}
		return fmt.Sprintf(`
type %s interface {
%s
}`, svc.GetName(), svcMethods(svc))
	}
}

func (svc *serviceDescription) GenerateRouterFunc(ctx TmplCtx, interfaceStyle string) (
	string, string) {
	sn := ctx.Get(VarContextServiceStructName)
	ctx.Set(VarContextServiceStructName, svc.GetName())
	a, b := GenTreeNodeRouteGen(ctx, svc, interfaceStyle)
	ctx.Set(VarContextServiceStructName, sn)
	return a, b
}

func (svc *serviceDescription) GenerateObjectString(objs []ObjTmpl, functions []FuncTmpl) string {
	var stringFragments []string

	for _, obj := range objs {
		stringFragments = append(stringFragments, obj.String())
	}

	for _, v := range functions {
		stringFragments = append(stringFragments, v.String())
	}

	return strings.Join(stringFragments, "\n")
}

func svcMethods(svc ServiceDescription) (res string) {
	res = fmt.Sprintf("    %sSignatureXXX() interface{}\n", svc.GetName())
	for _, cat := range svc.GetCategories() {
		res += _svcMethods(cat)
	}
	return
}

func _svcMethods(svc CategoryDescription) (res string) {
	for _, cat := range svc.GetCategories() {
		res += _svcMethods(cat)
	}
	for _, method := range svc.GetMethods() {
		res += "    " + method.GetName() + "(c controller.MContext)\n"
	}
	return
}
