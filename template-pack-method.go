package artisan

import "fmt"

func packMethod(objTmpl ObjTmpl, ctx TmplCtx) []FuncTmpl {
	if objTmpl.GetType() == TmplTypeEq {
		// skip eq func
		return nil
	}
	dependentFunc := vMethod(objTmpl, ctx)[0].(*FuncTmplImpl)
	dependentFunc.FD.Name = packInnerFunc(objTmpl)
	dependentFunc.Fix = true

	tmpl := newFuncPackTmpl(objTmpl)
	tmpl.FD.Name = fmt.Sprintf("PackSerialize%s", objTmpl.GetName())
	tmpl.ParameterList = objTmpl.GetSources()
	tmpl.ReturnType = pureType{typeString: objTmpl.GetName()}
	return []FuncTmpl{dependentFunc, tmpl}
}

func packInnerFunc(tmpl ObjTmpl) string {
	return fmt.Sprintf("_packSerialize%s", tmpl.GetName())
}

type funcPackTmplImpl struct {
	BaseFuncTmplImpl
}

func newFuncPackTmpl(rObject ObjTmpl) *funcPackTmplImpl {
	return &funcPackTmplImpl{
		BaseFuncTmplImpl: NewBaseFuncTmpl(true, rObject),
	}
}

func (f funcPackTmplImpl) String() string {
	if len(f.ParameterList) == 0 {
		return fmt.Sprintf(`func %s() (pack []%s) {
	return
}`, f.FD, f.ReturnType)
	}
	// dependent fp serialize
	return fmt.Sprintf(`func %s(%s) (pack []%s) {
	for i := range %s {
		pack = append(pack, %s(%s))
	}
	return
}`, f.FD, f.printParameterList(),
		f.ReturnType, f.ParameterList[0].Name,
		packInnerFunc(f.RObject), f.innerFuncParams())
}

func (f funcPackTmplImpl) innerFuncParams() (res string) {
	var appended = make(map[string]bool)
	for _, param := range f.ParameterList {
		if param.Source.GetSource() != nil {
			if _, ok := appended[param.Name]; ok {
				continue
			}
			appended[param.Name] = true
		}
		if len(res) != 0 {
			res += ", "
		}
		res += param.Name + "[i]"
	}
	return res
}

func (f funcPackTmplImpl) printParameterList() (res string) {
	var appended = make(map[string]bool)
	for _, param := range f.ParameterList {
		if param.Source.GetSource() != nil {
			if _, ok := appended[param.Name]; ok {
				continue
			}
			appended[param.Name] = true
		}
		if len(res) != 0 {
			res += ", "
		}
		res += param.Name + " []" + param.TypeOf
	}
	return res
}
