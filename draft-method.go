package artisan

type method struct {
	methodType MethodType
	name       string
	requests   []SerializeObject
	replies    []SerializeObject
}

func (method method) GetName() string {
	return method.name
}

func (method method) GetMethodType() MethodType {
	return method.methodType
}

func (method method) GetRequestProtocols() []SerializeObject {
	return method.requests
}

func (method method) GetResponseProtocols() []SerializeObject {
	return method.replies
}

func newMethod(methodType MethodType) *method {
	return &method{methodType: methodType}
}

func (method *method) CreateMethodDescription(ctx *Context) *methodDescription {
	desc := new(methodDescription)
	ctx.method = method
	desc.methodType = method.methodType
	desc.name = method.name
	for _, request := range method.requests {
		desc.requests = append(desc.requests, request.CreateObjectDescription(ctx))
	}
	for _, reply := range method.replies {
		desc.replies = append(desc.replies, reply.CreateObjectDescription(ctx))
	}
	return desc
}
