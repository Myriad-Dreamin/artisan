package artisan_core

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"unicode"
)

func searchSource(params []*XParam, sp ParameterDescription) *XParam {
	for _, param := range params {
		if bytes.Equal(sp.Hash(), param.Source.Hash()) {
			return param
		}
	}
	panic("not found")
}

func search(params []*XParam, sp ParameterDescription) string {
	source, param := sp.GetSource(), searchSource(params, sp)
	if source == nil {
		return param.Name
	} else {
		return param.Name + "." + source.MemberName()
	}
}

func getElements(i interface{}) (reflect.Value, reflect.Type) {
	return getReflectElements(reflect.ValueOf(i))
}

func getReflectElements(v reflect.Value) (reflect.Value, reflect.Type) {
	t := v.Type()
	for t.Kind() == reflect.Ptr {
		v, t = v.Elem(), t.Elem()
	}
	return v, t
}

func getElementValue(i interface{}) reflect.Value {
	v, _ := getElements(i)
	return v
}

func getElementType(i interface{}) reflect.Type {
	_, t := getElements(i)
	return t
}

func getReflectElementType(v reflect.Value) reflect.Type {
	_, t := getReflectElements(v)
	return t
}

func getReflectTypeElementType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func fromSnakeToCamel(src string, big bool) string {
	if len(src) == 0 {
		return ""
	}
	var b = bytes.NewBuffer(make([]byte, 0, len(src)))
	for i := range src {
		if src[i] == '_' {
			big = true
		} else {
			if big {
				big = false
				b.WriteByte(byte(unicode.ToUpper(rune(src[i]))))
			} else {
				b.WriteByte(src[i])
			}
		}
	}
	return b.String()
}

func allBig(name string) bool {
	for i := range name {
		if unicode.IsLower(rune(name[i])) {
			return false
		}
	}
	return true
}

func fromBigCamelToSnake(name string) string {

	if len(name) == 0 {
		return ""
	}
	var b = bytes.NewBuffer(make([]byte, 0, len(name)))
	b.WriteByte(byte(unicode.ToLower(rune(name[0]))))
	name = name[1:]
	var small = false
	for i := range name {
		if unicode.IsUpper(rune(name[i])) {
			if small {
				b.WriteByte('_')
				small = false
			}
			b.WriteByte(byte(unicode.ToLower(rune(name[i]))))
		} else {
			small = true
			b.WriteByte(name[i])
		}
	}
	return b.String()
}

func fromSnakeToSmallCamel(src string) string {
	return fromSnakeToCamel(src, false)
}

func fromSnakeToBigCamel(src string) string {
	return fromSnakeToCamel(src, true)
}

func toSmallCamel(name string) string {
	if len(name) == 0 {
		return name
	} else {
		return string(unicode.ToLower(rune(name[0]))) + name[1:]
	}
}

func getServiceName(name string) string {
	if len(name) == 0 {
		return "<embed-service>"
	}
	return name
}

type caller struct {
	fn, file string
	int
}

func (c caller) String() string {
	return fmt.Sprintf(`<function %s, file %s, line %d>`, c.fn, c.file, c.int)
}

func getCaller(skip int) caller {
	pc, f, l, _ := runtime.Caller(skip + 2)
	return caller{
		fn:   runtime.FuncForPC(pc).Name(),
		file: f,
		int:  l,
	}
}
