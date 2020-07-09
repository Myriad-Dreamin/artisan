package artisan

type RequestObject struct {
	s SerializeObject
}

func (r RequestObject) GetName() string {
	return r.s.GetName()
}

func (r RequestObject) DefiningPosition() string {
	return r.s.DefiningPosition()
}

func (r RequestObject) CreateObjectDescription(ctx *Context) ObjectDescription {
	ctx.Set("obj_suf", "Request")
	return r.s.CreateObjectDescription(ctx)
}
