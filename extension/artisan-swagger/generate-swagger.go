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
	"strconv"
	"strings"
	"time"
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

func toSwaggerStyle(s string) (_ string, wilds []string) {
	x := strings.Split(s, "/")
	for i := range x {
		if len(x[i]) > 0 && x[i][0] == ':' {
			wilds = append(wilds, x[i][1:])
			x[i] = "{" + x[i][1:] + "}"
		}
	}

	return strings.Join(x, "/"), wilds
}

func generateMethodsAndObjects(doc *spec.Swagger, tag *spec.Tag, base string, c artisan_core.CategoryDescription) {
	_ = c.GetMeta()
	_ = c.GetName()
	_ = c.GetObjects()

	subBase := path.Join(base, c.GetBase())

	for _, subCat := range c.GetCategories() {
		generateMethodsAndObjects(doc, tag, subBase, subCat)
	}

	var wilds []string

	subBase, wilds = toSwaggerStyle(subBase)

	var pathItem spec.PathItem

	for _, method := range c.GetMethods() {
		generatedOperation := generateOperation(doc, subBase, wilds, method)

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
		SchemaProps: spec.SchemaProps{
			Ref: spec.MustCreateRef("#/definitions/" + schemaName),
		},
	}
}

var responseGeneric = createRefSchema("genericResponse")

func generateOperation(doc *spec.Swagger, base string, wilds []string, method artisan_core.MethodDescription) *spec.Operation {
	var o = new(spec.Operation)

	o.ID = method.GetName()

	// generate path variables

	var param spec.Parameter
	param.In = "path"
	for i := range wilds {
		param.Name = wilds[i]
		o.Parameters = append(o.Parameters, param)
	}

	var inType string

	if method.GetMethodType() == artisan_core.GET {
		o.Consumes = []string{}
		o.Produces = []string{"application/json"}
		inType = "query"
	} else {
		o.Consumes = []string{"application/json", "application/x-www-form-urlencoded"}
		o.Produces = []string{"application/json"}
		inType = "body"
	}

	// generate requests

	requests := method.GetRequests()

	for i, request := range requests {

		if method.GetMethodType() == artisan_core.GET {
			for _, field := range request.GenObjectTmpl().GetFields() {
				var param spec.Parameter
				fieldDTOName := jsonFieldName(field.GetTag()["json"])
				if fieldDTOName == "-" {
					continue
				}
				requestSchema := createRuntimeProp(doc, field.GetType().(reflect.Type))
				param.Name = fieldDTOName
				param.Schema = &requestSchema

				param.In = inType
				o.Parameters = append(o.Parameters, param)
			}
		} else {
			var param spec.Parameter

			requestSchema := generateSchema(doc, request.GenObjectTmpl())
			param.Name = "request.option" + strconv.Itoa(i)
			param.Schema = &requestSchema

			param.In = inType
			o.Parameters = append(o.Parameters, param)
		}
	}

	// generate responses

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
	for _, reply := range method.GetReplies() {
		mainSchema.AllOf = append(mainSchema.AllOf, generateSchema(doc, reply.GenObjectTmpl()))
	}
	return o
}

func generateSchema(doc *spec.Swagger, transferObject artisan_core.ObjTmpl) spec.Schema {
	name := transferObject.GetName()
	if _, ok := doc.Definitions[name]; ok {
		return createRefSchema(name)
	}

	var schema spec.Schema

	schema.Type = []string{"object"}
	schema.Properties = make(map[string]spec.Schema)
	if doc.Definitions == nil {
		doc.Definitions = make(map[string]spec.Schema)
	}

	for _, field := range transferObject.GetFields() {
		fieldDTOName := jsonFieldName(field.GetTag()["json"])
		if fieldDTOName == "-" {
			continue
		}
		schema.Properties[fieldDTOName] = generateTypeSchema(doc, field.GetType())
	}

	doc.Definitions[name] = schema

	return createRefSchema(name)
}

func generateTypeSchema(doc *spec.Swagger, field artisan_core.Type) spec.Schema {
	switch ft := field.(type) {
	case reflect.Type:
		return createRuntimeProp(doc, ft)
	case artisan_core.ObjectDescType:
		return generateSchema(doc, ft.GenObjectTmpl())
	case artisan_core.ArrayType:
		var schema spec.Schema
		schema.Type = []string{"array"}
		schema.Items = new(spec.SchemaOrArray)
		var subSchema = generateTypeSchema(doc, ft.Type)
		schema.Items.Schema = &subSchema
		return schema
	default:
		panic("todo")
	}
}

func jsonFieldName(s string) string {
	return s
}

var timeType = reflect.TypeOf(new(time.Time))
var timeTypeElem = timeType.Elem()

func createRuntimeProp(doc *spec.Swagger, t reflect.Type) (schema spec.Schema) {
	if t == timeType || t == timeTypeElem {
		schema.Type = []string{"string"}
		schema.Format = "date-time"
		return
	}

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
		var xSchema = createRuntimeProp(doc, t.Elem())
		schema.Items.Schema = &xSchema
	case reflect.Array:
		schema.Type = []string{"array"}
		var xSchema = createRuntimeProp(doc, t.Elem())
		schema.Items.Schema = &xSchema
	case reflect.String:
		schema.Type = []string{"string"}
	case reflect.Float32:
		schema.Type = []string{"float"}
		schema.Format = "float32"
	case reflect.Float64:
		schema.Type = []string{"float"}
		schema.Format = "float64"
	case reflect.Ptr:
		return createRuntimeProp(doc, t.Elem())
	case reflect.Struct:
		var name = t.String()
		if _, ok := doc.Definitions[name]; ok {
			return
		}
		var externalSchema spec.Schema

		externalSchema.Type = []string{"object"}
		externalSchema.Properties = make(map[string]spec.Schema)
		if doc.Definitions == nil {
			doc.Definitions = make(map[string]spec.Schema)
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldDTOName := jsonFieldName(field.Tag.Get("json"))
			if fieldDTOName == "-" {
				continue
			}
			externalSchema.Properties[fieldDTOName] = createRuntimeProp(doc, field.Type)
		}

		doc.Definitions[name] = externalSchema
		return createRefSchema(name)
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
