package artisan

type ReplyObject struct {
	s SerializeObject
}

func (r ReplyObject) CreateObjectDescription(ctx *Context) ObjectDescription {
	ctx.set("obj_suf", "Reply")
	return r.s.CreateObjectDescription(ctx)
}
