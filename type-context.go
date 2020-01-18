package artisan

import (
	"reflect"
)

type Context struct {
	vars   map[string]interface{}
	method Method
	rawSvc ProposingService
	svc    ServiceDescription

	sources  map[uintptr]*source
	packages map[string]bool
}

func (c *Context) GetMethod() Method {
	return c.method
}

func (c *Context) GetRawSvc() ProposingService {
	return c.rawSvc
}

func (c *Context) GetSvc() ServiceDescription {
	return c.svc
}

func (c *Context) Clone() *Context {
	return &Context{
		vars:     c.vars,
		method:   c.method,
		rawSvc:   c.rawSvc,
		sources:  c.sources,
		packages: clonePackage(c.packages),
	}
}

func (c *Context) Sub() *Context {
	return &Context{
		vars:    c.vars,
		method:  c.method,
		rawSvc:  c.rawSvc,
		sources: c.sources,
	}
}

func (c *Context) AppendPackage(pkg string) {
	if len(pkg) != 0 {
		if c.packages == nil {
			c.packages = make(map[string]bool)
		}
		c.packages[pkg] = true
	}
}

func (c *Context) Set(k string, v interface{}) {
	if c.vars == nil {
		c.vars = make(map[string]interface{})
	}
	c.vars[k] = v
}

func (c *Context) Get(k string) (v interface{}) {
	if c.vars != nil {
		v, _ = c.vars[k]
	}
	return
}

func (c *Context) GetSource(ptr uintptr) *source {
	s, _ := c.sources[ptr]
	return s
}

func (c *Context) makeSources() {
	c.sources = make(map[uintptr]*source)
	models := c.rawSvc.GetModels()
	for _, xmodel := range models {
		v, t := reflect.ValueOf(xmodel.refer).Elem(), reflect.TypeOf(xmodel.refer).Elem()
		tt := t
		for t.Kind() == reflect.Ptr {
			v, t = v.Elem(), t.Elem()
		}
		if t.Kind() != reflect.Struct {
			panic(ErrNotStruct)
		}
		c.AppendPackage(t.PkgPath())
		for i := 0; i < t.NumField(); i++ {
			c.sources[v.Addr().Pointer()+t.Field(i).Offset] = &source{
				modelName: xmodel.name, faz: tt, fazElem: t, fieldIndex: i}
		}
	}
	//fmt.Println(c.sources)
}
