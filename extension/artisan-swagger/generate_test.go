package artisan_swagger

import (
	complex_example "github.com/Myriad-Dreamin/artisan/extension/artisan-swagger/example/complex"
	"testing"
)

func TestGenerateSwagger(t *testing.T) {
	GenerateSwagger(complex_example.Generate())
}
