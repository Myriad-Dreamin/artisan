package artisan

import "fmt"

func genStructFields(descriptions []ParameterDescription, xps []*XParam) (fields []ObjTmplField) {
	for i := range descriptions {
		desc := descriptions[i]
		fields = append(fields, ObjectTmplFieldImpl{
			Name:   desc.GetField().String(),
			PType:  desc.GetType(),
			Tag:    desc.GetTag(),
			Source: searchSource(xps, desc),
		})
	}
	return
}

func genTag(tags map[string]string) (res string) {
	res = "`"
	for k, v := range tags {
		if len(res) != 1 {
			res += " "
		}
		res += fmt.Sprintf(`%s:"%s"`, k, v)
	}
	res += "`"
	return res
}

func instantiateStructFields(fields []ObjTmplField) (res string) {
	for i := range fields {
		field := fields[i]
		if len(res) != 0 {
			res += "\n"
		}
		res += "    " + field.GetName() + " " + field.GetType().String() + " " + genTag(field.GetTag())
	}
	return
}
