package artisan_swagger

import (
	"encoding/json"
	"github.com/Myriad-Dreamin/artisan/artisan-core"
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
	"github.com/go-openapi/spec"
	"io/ioutil"
	"net/http"
	"path"
	"reflect"
	"strings"
)

func unmarshalHumanInfo(source interface{}, target interface{}) {
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

	if doc.Definitions != nil {
		var schema spec.Schema

		schema.Type = []string{"object"}
		schema.Properties = make(map[string]spec.Schema)

		var codeSchema spec.Schema
		codeSchema.Type = []string{"integer"}
		codeSchema.Format = "int"

		var errorSchema spec.Schema
		errorSchema.Type = []string{"string"}

		var paramSchema spec.Schema
		errorSchema.Type = []string{"object"}

		var paramsSchema spec.Schema
		paramsSchema.Type = []string{"object"}
		paramsSchema.Items = new(spec.SchemaOrArray)
		paramsSchema.Items.Schema = &paramSchema

		schema.Properties["code"] = codeSchema
		schema.Properties["error"] = errorSchema
		schema.Properties["params"] = paramsSchema

		doc.Definitions["genericResponse"] = schema
	}

	var toMerge = new(spec.Swagger)
	unmarshalHumanInfo(desc.HumanInfo, toMerge)
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
		generatedOperation := generateOperation(doc ,method)

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

func createRefSchema(schemaName string) spec.Schema {
	return spec.Schema{
		SchemaProps:        spec.SchemaProps{
			Ref: spec.MustCreateRef("#/definitions/" + schemaName),
		},
	}
}

var responseGeneric = createRefSchema("genericResponse")

func generateOperation(doc *spec.Swagger, method artisan_core.MethodDescription) *spec.Operation {
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
	resp.Description = "the request is serving. " +
		"however, the backend crashed and no response is returned."
	o.Responses.StatusCodeResponses[http.StatusInternalServerError] = resp
	resp = spec.Response{}
	resp.Description = "the request is served and a response is returned."
	var mainSchema = new(spec.Schema)
	resp.Schema = mainSchema
	o.Responses.StatusCodeResponses[http.StatusOK] = resp
	mainSchema.AllOf = append(mainSchema.AllOf, responseGeneric)
	for _, reply := range replies {
		generateSchema(doc, reply.GenObjectTmpl())
		mainSchema.AllOf = append(mainSchema.AllOf, createRefSchema(reply.GetType().String()))
	}
	return o
}

func generateSchema(doc *spec.Swagger, reply artisan_core.ObjTmpl) {
	var schema spec.Schema

	schema.Type = []string{"object"}
	schema.Properties = make(map[string]spec.Schema)
	if doc.Definitions == nil {
		doc.Definitions = make(map[string]spec.Schema)
	}
	if _, ok := doc.Definitions[reply.GetName()]; ok {
		return
	}

	for _, field := range reply.GetFields() {
		switch ft := field.GetType().(type) {
		case reflect.Type:
			schema.Properties[field.GetTag()["form"]] = createRuntimeProp(ft)
		case artisan_core.ObjectDescType:
			generateSchema(doc, ft.GenObjectTmpl())
			schema.Properties[field.GetTag()["form"]] = createRefSchema(ft.GetType().String())
		}
	}

	doc.Definitions[reply.GetName()] = schema
	return
}

func createRuntimeProp(t reflect.Type) (schema spec.Schema) {
	switch t.Kind() {
	case reflect.Int64:
		schema.Type = []string{"integer"}
		schema.Format = "int64"
	case reflect.Int32:
		schema.Type = []string{"integer"}
		schema.Format = "int32"
	case reflect.Int16:
		schema.Type = []string{"integer"}
		schema.Format = "int16"
	case reflect.Int8:
		schema.Type = []string{"integer"}
		schema.Format = "int8"
	case reflect.Int:
		schema.Type = []string{"integer"}
		schema.Format = "int"
	case reflect.Uint64:
		schema.Type = []string{"integer"}
		schema.Format = "uint64"
	case reflect.Uint32:
		schema.Type = []string{"integer"}
		schema.Format = "uint32"
	case reflect.Uint16:
		schema.Type = []string{"integer"}
		schema.Format = "uint16"
	case reflect.Uint8:
		schema.Type = []string{"integer"}
		schema.Format = "uint8"
	case reflect.Uint:
		schema.Type = []string{"integer"}
		schema.Format = "uint"
	case reflect.Bool:
		schema.Type = []string{"boolean"}
	case reflect.Slice:
		schema.Type = []string{"array"}
		schema.Items = new(spec.SchemaOrArray)
		var xSchema = createRuntimeProp(t.Elem())
		schema.Items.Schema = &xSchema
	case reflect.Array:
		schema.Type = []string{"array"}
		var xSchema = createRuntimeProp(t.Elem())
		schema.Items.Schema = &xSchema
	case reflect.String:
		schema.Type = []string{"string"}
	case reflect.Float32:
		schema.Type = []string{"float"}
		schema.Format = "float32"
	case reflect.Float64:
		schema.Type = []string{"float"}
		schema.Format = "float64"
	}
	return
}

func generateTag(v artisan_core.ServiceDescription) *spec.Tag {
	var tag spec.Tag

	tag.Name = v.GetName()

	// filePath := v.GetFilePath()
	//meta = v.GetMeta()
	//_ = meta

	unmarshalHumanInfo(v.GetHumanInfo(), &tag)
	return &tag
}
