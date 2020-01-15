package artisan

type ObjectTmplFieldImpl struct {
	Name   string
	PType  Type
	Tag    map[string]string
	Source *XParam
}

func (o ObjectTmplFieldImpl) GetName() string {
	return o.Name
}

func (o ObjectTmplFieldImpl) GetType() Type {
	return o.PType
}

func (o ObjectTmplFieldImpl) GetTag() map[string]string {
	return o.Tag
}

func (o ObjectTmplFieldImpl) GetSource() *XParam {
	return o.Source
}

func (o ObjectTmplFieldImpl) SetSource(s *XParam) {
	o.Source = s
	return
}
