package artisan

type RequestObject struct {
	s SerializeObject
}

func (r RequestObject) DefiningPosition() string {
	return r.s.DefiningPosition()
}

func (r RequestObject) CreateObjectDescription(ctx *Context) ObjectDescription {
	ctx.set("obj_suf", "Request")
	return r.s.CreateObjectDescription(ctx)
}
