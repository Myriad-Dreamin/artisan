package artisan

type serializeObject struct {
	dp     string
	uuid   UUID
	params []Parameter
	name   string
}

func newSerializeObject(skip int, descriptions ...interface{}) SerializeObject {
	var parameters []Parameter
	var name string
	for i := range descriptions {
		switch desc := descriptions[i].(type) {
		case SerializeObject:
			return desc
		case Parameter:
			parameters = append(parameters, desc)
		case []Parameter:
			parameters = append(parameters, desc...)
		case string:
			name = desc
		}
	}
	return &serializeObject{
		name:   name,
		params: parameters,
		uuid:   MustUUID(),
		dp:     getCaller(skip).String(),
	}
}

func (obj *serializeObject) DefiningPosition() string {
	return obj.dp
}

func (obj *serializeObject) GetName() string {
	return obj.name
}

func (obj serializeObject) GetUUID() UUID {
	return obj.uuid
}

func (obj *serializeObject) CreateObjectDescription(ctx *Context) ObjectDescription {
	desc := newObjectDescription(obj.uuid)
	for _, param := range obj.params {
		desc.params = append(desc.params, param.CreateParameterDescription(ctx))
	}
	desc.name = obj.name
	desc.dp = obj.dp
	if len(desc.name) == 0 {
		if suf := ctx.Get("obj_suf"); suf != nil {
			if suf, ok := suf.(string); ok {
				desc.name = ctx.method.GetName() + suf
			}
		} else {
			panic(errObjectHasNoName(obj, ctx))
		}
	}

	//fmt.Println("creating", desc.name)
	//for i := range desc.params {
	//	fmt.Print("    ")
	//	param := desc.params[i]
	//	fmt.Println(param.fieldName, param.typeString, param.tags, "||", param.embedObjects)
	//}
	return desc
}
