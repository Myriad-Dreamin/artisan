package artisan_core

import "fmt"

func pMethod(objTmpl ObjTmpl, _ TmplCtx) []FuncTmpl {
	if objTmpl.GetType() == TmplTypeEq {
		// skip eq func
		return nil
	}
	tmpl := NewFuncTmpl(false, objTmpl)
	tmpl.FD.Name = fmt.Sprintf("PSerialize%s", objTmpl.GetName())
	tmpl.ParameterList = objTmpl.GetSources()
	tmpl.ReturnType = pureType{typeString: fmt.Sprintf("*%s", objTmpl.GetName())}
	return []FuncTmpl{tmpl}
}
