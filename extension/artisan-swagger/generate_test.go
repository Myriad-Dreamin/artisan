package artisan_swagger

import (
	"encoding/json"
	"fmt"
	complex_example "github.com/Myriad-Dreamin/artisan/extension/artisan-swagger/example/complex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateSwagger(t *testing.T) {
	b, err := json.Marshal(GenerateSwagger(complex_example.Generate()))
	assert.NoError(t, err)
	fmt.Println(string(b))
}
