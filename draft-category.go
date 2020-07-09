package artisan

type category struct {
	name        string
	path        string
	methods     []Method
	wildObjects []SerializeObject
	subs        map[string]Category
}

func newCategory() *category {
	return new(category)
}

func (c *category) GetName() string {
	return c.name
}

func (c *category) GetPath() string {
	return c.path
}

func (c *category) GetMethods() []Method {
	return c.methods
}

func (c *category) GetWildObjects() []SerializeObject {
	return c.wildObjects
}

func (c *category) ForEachSubCate(mapFunc func(path string, cat Category) (shouldStop bool)) error {
	for k, v := range c.subs {
		if !mapFunc(k, v) {
			return ErrorStopped
		}
	}
	return nil
}

func (c *category) WithName(name string) Category {
	c.name = name
	return c
}

func (c *category) Path(path string) Category {
	c.path = path
	return c
}

func (c *category) SubCate(path string, cat Category) Category {
	if _, ok := c.subs[path]; ok {
		panic(ErrConflictPath)
	}
	c.subs[path] = cat
	return c
}

func (c *category) DiveIn(path string) Category {
	cat := &category{
		path: path,
	}
	c.SubCate(path, cat)
	return cat
}

func (c *category) RawMethod(m ...Method) Category {
	c.methods = append(c.methods, m...)
	return c
}

// todo
func (c *category) Method(m MethodType, descriptions ...interface{}) Category {
	method := newMethod(m)
	for _, description := range descriptions {
		switch desc := description.(type) {
		case string:
			method.name = desc
		case RequestObject:
			method.requests = append(method.requests, desc)
		case ReplyObject:
			method.replies = append(method.replies, desc)
		}
	}

	c.methods = append(c.methods, method)
	return c
}

func (c *category) Object(descriptions ...interface{}) Category {
	c.wildObjects = append(c.wildObjects, newSerializeObject(1, descriptions...))
	return c
}

func (c *category) HelpWrapObjectXXX(skip int, descriptions ...interface{}) Category {
	c.wildObjects = append(c.wildObjects, newSerializeObject(skip+1, descriptions...))
	return c
}

func (c *category) AppendObject(objs ...SerializeObject) Category {
	c.wildObjects = append(c.wildObjects, objs...)
	return c
}

func (c *category) CreateCategoryDescription(ctx *Context) CategoryDescription {
	desc := new(categoryDescription)
	for _, method := range c.methods {
		desc.methods = append(desc.methods, method.CreateMethodDescription(ctx.Sub()))
	}
	for _, obj := range c.wildObjects {
		desc.objDesc = append(desc.objDesc, obj.CreateObjectDescription(ctx))
	}

	for _, sub := range c.subs {
		subDesc := sub.CreateCategoryDescription(ctx.Clone())
		desc.subCates[subDesc.GetName()] = subDesc
	}
	desc.name = c.name
	desc.path = c.path
	return desc
}
