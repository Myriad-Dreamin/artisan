package artisan_core

type objectDescription struct {
	name   string
	dp     string
	uuid   UUID
	params []ParameterDescription

	cache ObjTmpl
}

func newObjectDescription(uuid UUID) *objectDescription {
	return &objectDescription{uuid: uuid}
}

func (desc *objectDescription) GetName() string {
	return desc.name
}

func (desc *objectDescription) DefiningPosition() string {
	return desc.dp
}

func (desc *objectDescription) GetUUID() UUID {
	return desc.uuid
}

func (desc *objectDescription) GetContainingParams() []ParameterDescription {
	return desc.params
}

func (desc *objectDescription) GetEmbedObject() (dx []ObjectDescription) {
	for _, param := range desc.params {
		dx = append(dx, param.GetEmbedObjects()...)
	}
	return dx
}

func (desc *objectDescription) GetPackages() PackageSet {
	var pac PackageSet
	for _, param := range desc.params {
		pac = PackageSetInPlaceMerge(pac, param.GetPackages())
	}
	return pac
}

func (desc *objectDescription) GetType() Type {
	return ObjectDescType{desc}
}

func (desc *objectDescription) GenObjectTmpl() ObjTmpl {
	if desc.cache != nil {
		return desc.cache
	}

	xps := desc.genXParams()
	desc.cache = &ObjTmplImpl{
		// type desc.name struct {
		Name: desc.name, TType: TmplTypeStruct,
		Fields: genStructFields(desc.params, xps),
		Xps:    xps,
	}
	return desc.cache
}

func (desc *objectDescription) genXParams() (params []*XParam) {
	//desc.params
	for _, param := range desc.params {
		source := param.GetSource()
		if source != nil {
			params = append(params, &XParam{
				Name:   source.ParamName(),
				TypeOf: source.faz.String(),
				Source: param,
			})
		} else {
			params = append(params, &XParam{
				Name:   "_" + fromSnakeToSmallCamel(param.GetDTOName()),
				TypeOf: param.GetType().String(),
				Source: param,
			})
		}
	}
	return
}
