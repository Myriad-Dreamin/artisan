package artisan

import "fmt"

func vMethod(objTmpl ObjTmpl, _ TmplCtx) []FuncTmpl {
	if objTmpl.GetType() == TmplTypeEq {
		// skip eq func
		return nil
	}
	tmpl := NewFuncTmpl(false, objTmpl)
	tmpl.FD.Name = fmt.Sprintf("Serialize%s", objTmpl.GetName())
	tmpl.ParameterList = objTmpl.GetSources()
	tmpl.ReturnType = pureType{typeString: objTmpl.GetName()}
	return []FuncTmpl{tmpl}
}
