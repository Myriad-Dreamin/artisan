package artisan

type categoryDescription struct {
	name          string
	path          string
	meta          interface{}
	subCates      map[string]CategoryDescription
	tmplFactories []FuncTmplFac
	methods       []MethodDescription
	objDesc       []ObjectDescription
	packages      PackageSet
}

func (c *categoryDescription) GetPath() string {
	return c.path
}

func (c *categoryDescription) GetBase() string {
	return c.path
}

func (c *categoryDescription) GetMeta() interface{} {
	return c.meta
}

func (c *categoryDescription) GetName() string {
	return c.name
}

func (c *categoryDescription) GetMethods() []MethodDescription {
	return c.methods
}

func (c *categoryDescription) GetPackages() PackageSet {
	if c == nil {
		return nil
	}
	pac := make(PackageSet)
	for _, method := range c.methods {
		pac = PackageSetInPlaceMerge(pac, method.GetPackages())
	}
	for _, obj := range c.objDesc {
		pac = PackageSetInPlaceMerge(pac, obj.GetPackages())
	}
	//for _, cate := range c.subCates {
	//	pac = PackageSetInPlaceMerge(pac, cate.GetPackages())
	//}
	return pac
}

func (c *categoryDescription) IterCategories(callback func(k string, v CategoryDescription) bool) bool {
	if c.subCates != nil {
		for k, v := range c.subCates {
			if callback(k, v) {
				return true
			}
		}
	}
	return false
}

func (c *categoryDescription) GetCategories() (categories []CategoryDescription) {
	for _, x := range c.subCates {
		categories = append(categories, x)
	}
	return
}

func (c *categoryDescription) GetObjects() []ObjectDescription {
	return c.objDesc
}

func (c *categoryDescription) GetTmplFunctionFactory() []FuncTmplFac {
	return c.tmplFactories
}

func (c *categoryDescription) SetName(n string) CategoryDescription {
	c.name = n
	return c
}

func (c *categoryDescription) GenerateRouterFunc(ctx TmplCtx, interfaceStyle string) (
	string, string) {
	return GenTreeNodeRouteGen(ctx, c, interfaceStyle)
}

func (c *categoryDescription) GenerateObjects(ts []FuncTmplFac, ctx TmplCtx) (objs []ObjTmpl, funcs []FuncTmpl) {
	return GenerateObjects(c, ts, ctx)
}
