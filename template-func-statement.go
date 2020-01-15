package artisan

// todo statement
type Statement struct {
	Dst      []*XParam
	Src      []*XParam
	Caller   FuncDescription
	HasError bool
}
