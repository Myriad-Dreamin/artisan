package artisan_core

type arrayParam struct {
	p Parameter
}

type arrayParamDescription struct {
	ParameterDescription
	a ArrayType
}

type ArrayType struct {
	Type
	typeString string
}

func (a ArrayType) String() string {
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
	desc.a = ArrayType{Type: desc.ParameterDescription.GetType()}
	return desc
}
