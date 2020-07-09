package artisan

import (
	"errors"
	"reflect"
)

type norm struct {
	name  string
	tags  []*tag
	field Field
	param reflect.Value
}

func newNorm(name string) *norm {
	return &norm{name: name}
}

type FieldName string

func createNorm(name string, descriptions ...interface{}) *norm {
	param := newNorm(name)
	for _, description := range descriptions {
		switch desc := description.(type) {
		case *tag:
			param.tags = append(param.tags, desc)
		case FieldName:
			param.field = pureField{string(desc)}
		case SerializeObject:
			param.param = reflect.ValueOf(desc)
		default:
			param.param = reflect.ValueOf(desc).Elem()
		}
	}
	return param
}

func (n *norm) CreateParameterDescription(ctx *Context) ParameterDescription {
	desc := new(parameterDescription)
	desc.name = n.name
	desc.field = n.field
	if desc.field == nil {
		desc.field = pureField{fromSnakeToBigCamel(desc.name)}
	}

	if embedType, ok := n.param.Interface().(SerializeObject); ok {
		objDesc := embedType.CreateObjectDescription(ctx)
		desc.embedObjects = append(desc.embedObjects, objDesc)
		desc.pType = objDesc.GetType()
	} else {
		desc.calculatedPackages = PackageSetAppend(desc.calculatedPackages,
			getReflectElementType(n.param).PkgPath())
		desc.pType = parseParamType(desc.calculatedPackages, n)
		desc.source = parseSource(ctx, n)
		if desc.source != nil {
			desc.calculatedPackages = PackageSetInPlaceMerge(desc.calculatedPackages,
				desc.source.GetPackageSet())
		}
	}
	desc.tags = make(map[string]string)
	desc.tags["json"] = desc.name
	desc.tags["form"] = desc.name
	for _, tag := range n.tags {
		if v, ok := desc.tags[tag.key]; ok {
			desc.tags[tag.key] = v + ";" + tag.value
		} else {
			desc.tags[tag.key] = tag.value
		}
	}

	return desc
}

func parseSource(context *Context, n *norm) *source {
	return context.GetSource(n.param.UnsafeAddr())
}

func parseParamType(packageSet PackageSet, n *norm) Type {
	t := n.param.Type()
	if t != nil {
		if packageSet == nil {
			// todo: norm pos
			panic(errors.New("package set nil"))
		}
		PackageSetAppend(packageSet, t.PkgPath())
		return t
	} else {
		panic("nil type")
	}
}
