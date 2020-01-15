package artisan

type model struct {
	name  string
	refer interface{}
}

type Name string

func Model(descriptions ...interface{}) *model {
	m := new(model)
	for i := range descriptions {
		switch desc := descriptions[i].(type) {
		case Name:
			m.name = string(desc)
		default:
			m.refer = desc
		}
	}
	return m
}
