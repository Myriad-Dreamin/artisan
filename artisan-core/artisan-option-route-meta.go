package artisan_core

type RouterMeta struct {
	RuntimeRouterMeta string
	NeedAuth          bool
}

func (r RouterMeta) GetRuntimeRouterMeta() string {
	return r.RuntimeRouterMeta
}

func (r RouterMeta) GetNeedAuth() bool {
	return r.NeedAuth
}
