package artisan

import "fmt"

type FuncTmplImpl struct {
	BaseFuncTmplImpl
}

func (f FuncTmplImpl) String() string {
	return fmt.Sprintf(`func %s(%s) %s {
%s
    return %s
}`, f.FD, f.printParamsList(), f.ReturnType, f.printStatments(),
		f.genResultStruct())
}

func (f FuncTmplImpl) printParamsList() (res string) {
	p := f.ParameterList
	var appended = make(map[string]bool)
	for _, param := range p {
		if param.Source != nil {
			if _, ok := appended[param.Name]; ok {
				continue
			}
			appended[param.Name] = true
		}
		if len(res) != 0 {
			res += ", "
		}
		res += param.Name + " " + param.TypeOf
	}
	return res
}

func (f FuncTmplImpl) genResultStruct() (res string) {
	for _, field := range f.RObject.GetFields() {
		if len(res) != 0 {
			res += "\n"
		}
		res += "        " + field.GetName() + ": " + XParamAsVar(field.GetSource()) + ","
	}
	res = fmt.Sprintf(`%s{
%s
    }`, getCreator(f.ReturnType.String()), res)
	return
}

func getCreator(s string) string {
	if len(s) == 0 {
		panic("empty type string")
	}
	if s[0] == '*' {
		if len(s) == 1 {
			panic("invalid type string")
		}
		if s[1] == '*' {
			panic("invalid ptr of ptr")
		}
		return "&" + s[1:]
	} else {
		return s
	}
}
