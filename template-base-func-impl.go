package artisan

type BaseFuncTmplImpl struct {
	FD      FuncDescription
	RObject ObjTmpl

	Fix           bool
	CanError      bool
	ParameterList []*XParam
	ReturnType    Type
	Statements    []*Statement
}

func (f BaseFuncTmplImpl) WantFix() bool {
	return f.Fix
}

func (f BaseFuncTmplImpl) String() string {
	return "//base func impl does not generate anything"
}

func (f BaseFuncTmplImpl) AllowError() bool {
	return f.CanError
}

func (f BaseFuncTmplImpl) GetResponsibleObject() ObjTmpl {
	return f.RObject
}

func (f BaseFuncTmplImpl) GetName() string {
	return f.FD.Name
}

func (f BaseFuncTmplImpl) GetReceiver() *FuncTmplReceiver {
	return f.FD.Receiver
}

func (f BaseFuncTmplImpl) GetParameterList() []*XParam {
	return f.ParameterList
}

func (f BaseFuncTmplImpl) GetReturnType() Type {
	return f.ReturnType
}

func (f BaseFuncTmplImpl) GetStatements() []*Statement {
	return f.Statements
}

func (f BaseFuncTmplImpl) printStatments() string {
	if len(f.Statements) != 0 {
		return `    //todo statements`
	}
	return ""
}
