package artisan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInnerFromBigCamelToSnake(t *testing.T) {
	assert.Equal(t, fromBigCamelToSnake("test"), "test")
	assert.Equal(t, fromBigCamelToSnake("Test"), "test")
	assert.Equal(t, fromBigCamelToSnake("TestID"), "test_id")
	assert.Equal(t, fromBigCamelToSnake("ABC"), "abc")
	assert.Equal(t, fromBigCamelToSnake("QwQ"), "qw_q")
}
