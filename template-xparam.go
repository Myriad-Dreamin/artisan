package artisan

type XParam struct {
	Name, TypeOf string
	Source       ParameterDescription
}

func XParamAsVar(param *XParam) string {
	src := param.Source
	if src != nil && src.GetSource() != nil {
		return param.Name + "." + src.GetSource().MemberName()
	} else {
		return param.Name
	}
}
