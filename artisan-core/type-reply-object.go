package artisan_core

type ReplyObject struct {
	s SerializeObject
}

func (r ReplyObject) GetName() string {
	return r.s.GetName()
}

func (r ReplyObject) DefiningPosition() string {
	return r.s.DefiningPosition()
}

func (r ReplyObject) CreateObjectDescription(ctx *Context) ObjectDescription {
	ctx.Set("obj_suf", "Reply")
	return r.s.CreateObjectDescription(ctx)
}
