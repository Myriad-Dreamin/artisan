package artisan

type serializeObject struct {
	params []Parameter
	name   string
}

func (obj *serializeObject) CreateObjectDescription(ctx *Context) ObjectDescription {
	desc := new(objectDescription)
	for _, param := range obj.params {
		desc.params = append(desc.params, param.CreateParameterDescription(ctx))
	}
	desc.name = obj.name
	if len(desc.name) == 0 {
		if suf := ctx.get("obj_suf"); suf != nil {
			if suf, ok := suf.(string); ok {
				desc.name = ctx.method.GetName() + suf
			}
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
