package artisan

import (
	"reflect"
	"time"
)

type pureType struct {
	typeString string
}

func (p pureType) String() string {
	return p.typeString
}

type pureField struct {
	s string
}

func (p pureField) String() string {
	return p.s
}

//noinspection GoUnusedGlobalVariable
var (
	_interface = interface{}(0)
	Interface  = &_interface
	Bool       = new(bool)
	String     = new(string)
	Strings    = new([]string)
	Byte       = new(byte)
	Bytes      = new([]byte)
	Rune       = new(rune)
	Runes      = new([]rune)
	Time       = new(time.Time)
	PTime      = new(*time.Time)

	Int   = new(int)
	Int8  = new(int8)
	Int16 = new(int16)
	Int32 = new(int32)
	Int64 = new(int64)

	Uint   = new(uint)
	Uint8  = new(uint8)
	Uint16 = new(uint16)
	Uint32 = new(uint32)
	Uint64 = new(uint64)
)

//noinspection GoUnusedGlobalVariable
var cateType = reflect.TypeOf(new(Category))
