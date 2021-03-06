package artisan_core

type methodDescription struct {
	methodType MethodType
	name       string
	authMeta   string
	requests   []ObjectDescription
	replies    []ObjectDescription
}

func (m methodDescription) GetMethodType() MethodType {
	return m.methodType
}

func (m methodDescription) GetName() string {
	return m.name
}

func (m methodDescription) GetAuthMeta() string {
	return m.authMeta
}

func (m methodDescription) GetRequests() []ObjectDescription {
	return m.requests
}

func (m methodDescription) GetReplies() []ObjectDescription {
	return m.replies
}

func (m methodDescription) GetPackages() PackageSet {
	var pac PackageSet
	for _, req := range m.requests {
		pac = PackageSetInPlaceMerge(pac, req.GetPackages())
	}
	for _, res := range m.replies {
		pac = PackageSetInPlaceMerge(pac, res.GetPackages())
	}
	return pac
}
