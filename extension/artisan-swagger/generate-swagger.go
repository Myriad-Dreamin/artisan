package artisan_swagger

import (
	"encoding/json"
	"github.com/Myriad-Dreamin/artisan/artisan-core"
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
	"github.com/go-openapi/spec"
	"io/ioutil"
)

func GenerateSwagger(desc *artisan_core.PublishedServices) *spec.Swagger {
	var doc = new(spec.Swagger)
	var err error
	if desc.HumanInfo != nil {
		var b []byte
		switch hi := desc.HumanInfo.(type) {
		case string:
			b, err = ioutil.ReadFile(hi)
			sugar.HandlerError0(err)
		default:
			b, err = json.Marshal(hi)
			sugar.HandlerError0(err)
			err = json.Unmarshal(b, doc)
			sugar.HandlerError0(err)
		}
	}

	return doc
}
