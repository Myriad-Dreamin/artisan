package artisan

func SetDefaultFunctionTmplFactories(funcs []FuncTmplFac) {
	DefaultFunctionTmplFactories = funcs
}

func SetDefaultMetaFactory(_metaFac func() interface{}) {
	metaFac = _metaFac
}

func SetDefaultNewTmplContext(_newTmplContext func(svc ServiceDescription) TmplCtx) {
	newTmplContext = _newTmplContext
}

var DefaultFunctionTmplFactories = []FuncTmplFac{
	pMethod, vMethod, packMethod,
}

var metaFac = defaultMeta

func defaultMeta() interface{} { return nil }

var newTmplContext = defaultNewTmplContext

func defaultNewTmplContext(svc ServiceDescription) TmplCtx {
	return &TmplContextImpl{
		packages: make(PackageSet),
		uuid:     make(map[string]bool),
		svc:      svc,
		meta:     metaFac(),
	}
}
