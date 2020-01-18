package artisan

import (
	"fmt"
	"reflect"
)

func GenerateObjects(
	g GenTreeNode, ts []FuncTmplFac, c TmplCtx) (
	objs []ObjTmpl, funcs []FuncTmpl) {
	ctx := c.Clone()
	ctx.MergePackages(g.GetPackages())
	tmplFactories := append(g.GetTmplFunctionFactory(), ts...)

	for _, cat := range g.GetCategories() {
		ctx.PushCategory(cat)

		os, fs := cat.GenerateObjects(tmplFactories, ctx)
		objs = append(objs, os...)
		funcs = append(funcs, fs...)

		for _, method := range cat.GetMethods() {
			ctx.SetObjectType(ObjectTypeRequest)
			for _, req := range method.GetRequests() {

				os, fs := dumpObj(ctx, tmplFactories, req)
				objs = append(objs, os...)
				funcs = append(funcs, fs...)
			}

			ctx.SetObjectType(ObjectTypeReply)
			for _, res := range method.GetReplies() {
				os, fs := dumpObj(ctx, tmplFactories, res)
				objs = append(objs, os...)
				funcs = append(funcs, fs...)
			}
		}

		ctx.SetObjectType(ObjectTypeObject)
		for _, obj := range cat.GetObjects() {
			os, fs := dumpObj(ctx, tmplFactories, obj)
			objs = append(objs, os...)
			funcs = append(funcs, fs...)
		}

		ctx.PopCategory()
	}
	return
}

func dumpObj(ctx TmplCtx, factories []FuncTmplFac,
	desc ObjectDescription) (objs []ObjTmpl, funcs []FuncTmpl) {
	fmt.Println("testing", reflect.TypeOf(desc), desc.GetUUID())
	if !ctx.AppendUUID(desc.GetUUID()) {
		return
	}
	fmt.Println(desc.GetType())

	tmpl := desc.GenObjectTmpl()
	objs = append(objs, tmpl)
	for _, fac := range factories {
		fs := fac(tmpl, ctx)
		for i := range fs {
			f := fs[i]
			if !f.WantFix() {
				//	todo middleware
			}
			fs[i] = f
		}
		funcs = append(funcs, fs...)
	}

	for _, obj := range desc.GetEmbedObject() {
		os, fs := dumpObj(ctx, factories, obj)
		objs = append(objs, os...)
		funcs = append(funcs, fs...)
	}
	return objs, funcs
}
