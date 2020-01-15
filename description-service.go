package artisan

type serviceDescription struct {
	name          string
	base          string
	tmplFactories []FuncTmplFac
	categories    []CategoryDescription
	filePath      string
	//packages   map[string]int
}

func (description serviceDescription) GetPackages() PackageSet {
	return nil
}

func (description serviceDescription) GetTmplFunctionFactory() []FuncTmplFac {
	return description.tmplFactories
}

func (description serviceDescription) GetName() string {
	return description.name
}

func (description serviceDescription) GetBase() string {
	return description.base
}

func (description serviceDescription) GetCategories() []CategoryDescription {
	return description.categories
}

func (description serviceDescription) GetFilePath() string {
	return description.filePath
}

func (description serviceDescription) GenerateObjects(ts []FuncTmplFac, c TmplCtx) (objs []ObjTmpl, funcs []FuncTmpl) {
	return GenerateObjects(description, ts, c)
}
