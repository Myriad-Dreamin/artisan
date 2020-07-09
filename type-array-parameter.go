package artisan

type arrayParam struct {
	p Parameter
}

type arrayParamDescription struct {
	ParameterDescription
	a arrayType
}

type arrayType struct {
	Type
	typeString string
}

func (a arrayType) String() string {
	return "[]" + a.Type.String()
}

func (a *arrayParamDescription) GetPackages() PackageSet {
	return a.ParameterDescription.GetPackages()
}

func (a *arrayParamDescription) GetType() Type {
	return a.a
}

func (a arrayParam) CreateParameterDescription(ctx *Context) ParameterDescription {
	desc := &arrayParamDescription{
		ParameterDescription: a.p.CreateParameterDescription(ctx),
	}
	desc.a = arrayType{Type: desc.ParameterDescription.GetType()}
	return desc
}
