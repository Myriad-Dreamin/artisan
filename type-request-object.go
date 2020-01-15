package artisan

type RequestObject struct {
	s SerializeObject
}

func (r RequestObject) CreateObjectDescription(ctx *Context) ObjectDescription {
	ctx.set("obj_suf", "Request")
	return r.s.CreateObjectDescription(ctx)
}
