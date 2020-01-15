package artisan

type methodDescription struct {
	methodType MethodType
	name       string
	requests   []ObjectDescription
	replies    []ObjectDescription
}

func (m methodDescription) GetMethodType() MethodType {
	return m.methodType
}

func (m methodDescription) GetName() string {
	return m.name
}

func (m methodDescription) GetRequests() []ObjectDescription {
	return m.requests
}

func (m methodDescription) GetReplies() []ObjectDescription {
	return m.replies
}
