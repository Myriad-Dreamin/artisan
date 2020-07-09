package artisan

type categoryDescription struct {
	name          string
	path          string
	subCates      map[string]CategoryDescription
	tmplFactories []FuncTmplFac
	methods       []MethodDescription
	objDesc       []ObjectDescription
	packages      PackageSet
}

func (c *categoryDescription) GetObjects() []ObjectDescription {
	return c.objDesc
}

func (c *categoryDescription) GenerateObjects(ts []FuncTmplFac, ctx TmplCtx) (objs []ObjTmpl, funcs []FuncTmpl) {
	return GenerateObjects(c, ts, ctx)
}

func (c *categoryDescription) GetTmplFunctionFactory() []FuncTmplFac {
	return c.tmplFactories
}

func (c *categoryDescription) GetPath() string {
	return c.path
}

func (c *categoryDescription) GetCategories() (categories []CategoryDescription) {
	for _, x := range c.subCates {
		categories = append(categories, x)
	}
	return
}

func (c *categoryDescription) GetMethods() []MethodDescription {
	return c.methods
}

func (c *categoryDescription) GetName() string {
	return c.name
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
