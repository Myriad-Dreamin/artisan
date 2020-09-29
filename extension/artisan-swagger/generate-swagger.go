package artisan_swagger

import (
	"encoding/json"
	"github.com/Myriad-Dreamin/artisan/artisan-core"
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
	"github.com/go-openapi/spec"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

func unwrapHumanInfo(source interface{}, target interface{}) {
	var err error

	if source != nil {
		var b []byte
		switch hi := source.(type) {
		case string:
			b, err = ioutil.ReadFile(hi)
			sugar.HandlerError0(err)
			err = json.Unmarshal(b, target)
			sugar.HandlerError0(err)
		default:
			b, err = json.Marshal(hi)
			sugar.HandlerError0(err)
			err = json.Unmarshal(b, target)
			sugar.HandlerError0(err)
		}
	}
}

type mergingSwagger spec.Swagger

func (s *mergingSwagger) Merge(toMerge *spec.Swagger) {
	s.Info = toMerge.Info
}

func GenerateSwagger(desc *artisan_core.PublishedServices) *spec.Swagger {
	var doc = new(spec.Swagger)
	doc.Swagger = "2.0"

	desc.GetPackageName()

	doc.BasePath = desc.Base

	doc.Tags = make([]spec.Tag, 0, len(desc.SvcMap))
	doc.Paths = new(spec.Paths)
	for _, v := range desc.SvcMap {
		tag := generateTag(v)
		doc.Tags = append(doc.Tags, *tag)

		base := v.GetBase()
		if !strings.HasPrefix(base, "/") {
			base = "/" + base
		}
		for _, c := range v.GetCategories() {
			generateMethodsAndObjects(doc, tag, base, c)
		}
	}

	var toMerge = new(spec.Swagger)
	unwrapHumanInfo(desc.HumanInfo, toMerge)
	(*mergingSwagger)(doc).Merge(toMerge)
	return doc
}

func generateMethodsAndObjects(doc *spec.Swagger, tag *spec.Tag, base string, c artisan_core.CategoryDescription) {
	_ = c.GetMeta()
	_ = c.GetName()
	_ = c.GetObjects()

	subBase := path.Join(base, c.GetBase())

	for _, subCat := range c.GetCategories() {
		generateMethodsAndObjects(doc, tag, subBase, subCat)
	}

	var pathItem spec.PathItem

	for _, method := range c.GetMethods() {
		generatedOperation := generateOperation(method)

		generatedOperation.Tags = []string{tag.Name}
		switch method.GetMethodType() {
		case artisan_core.POST:
			pathItem.Post = generatedOperation
		case artisan_core.GET:
			pathItem.Get = generatedOperation
		case artisan_core.PATCH:
			pathItem.Patch = generatedOperation
		case artisan_core.HEAD:
			pathItem.Head = generatedOperation
		case artisan_core.PUT:
			pathItem.Put = generatedOperation
		case artisan_core.DELETE:
			pathItem.Delete = generatedOperation
		default:
			panic("not supported...")
		}
	}

	if doc.Paths.Paths == nil {
		doc.Paths.Paths = make(map[string]spec.PathItem)
	}

	doc.Paths.Paths[subBase] = pathItem
}

func generateOperation(method artisan_core.MethodDescription) *spec.Operation {
	var o = new(spec.Operation)
	_ = method.GetRequests()
	replies := method.GetReplies()

	o.ID = method.GetName()
	if method.GetMethodType() == artisan_core.GET {
		o.Consumes = []string{"querystring"}
		o.Produces = []string{"application/json"}

	} else {
		o.Consumes = []string{"application/json", "application/x-www-form-urlencoded"}
		o.Produces = []string{"application/json"}
	}

	o.Responses = new(spec.Responses)
	o.Responses.StatusCodeResponses = make(map[int]spec.Response)

	var resp = spec.Response{}
	resp.Description = "the request is served and a response is returned."
	o.Responses.StatusCodeResponses[http.StatusOK] = resp
	resp = spec.Response{}
	resp.Description = "the request is serving." +
		"however, the backend fall back and no response is returned."
	o.Responses.StatusCodeResponses[http.StatusInternalServerError] = resp

	switch len(replies) {
	case 0:
	case 1:
	default:
	}
	return o
}

func generateTag(v artisan_core.ServiceDescription) *spec.Tag {
	var tag spec.Tag

	tag.Name = v.GetName()

	// filePath := v.GetFilePath()
	//meta = v.GetMeta()
	//_ = meta

	unwrapHumanInfo(v.GetHumanInfo(), &tag)
	return &tag
}
