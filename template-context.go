package artisan

type TmplContextImpl struct {
	packages PackageSet

	objType    ObjectType
	svc        ServiceDescription
	categories []CategoryDescription

	meta interface{}
}

func (t *TmplContextImpl) GetPackages() PackageSet {
	return t.packages
}

func (t *TmplContextImpl) MergePackages(pks PackageSet) {
	t.packages = inplaceMergePackage(t.packages, pks)
}

func (t *TmplContextImpl) SetObjectType(oT ObjectType) {
	t.objType = oT
}

func (t *TmplContextImpl) CurrentObjectType() ObjectType {
	return t.objType
}

func (t *TmplContextImpl) Clone() TmplCtx {
	return &TmplContextImpl{
		packages:   t.packages,
		svc:        t.svc,
		categories: t.categories,
		meta:       t.meta,
	}
}

func (t *TmplContextImpl) PushCategory(cat CategoryDescription) {
	t.categories = append(t.categories, cat)
	return
}

func (t *TmplContextImpl) PopCategory() (cat CategoryDescription) {
	t.categories, cat = t.categories[:len(t.categories)-1], t.categories[len(t.categories)-1]
	return
}

func (t *TmplContextImpl) AppendPackage(pkgPath string) {
	t.packages[pkgPath] = true
}

func (t *TmplContextImpl) GetService() ServiceDescription {
	return t.svc
}

func (t *TmplContextImpl) GetCategories() []CategoryDescription {
	return t.categories
}

func (t *TmplContextImpl) SetMeta(meta interface{}) {
	t.meta = meta
}

func (t *TmplContextImpl) GetMeta() interface{} {
	return t.meta
}
