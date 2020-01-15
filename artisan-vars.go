package artisan

import "errors"

var (
	ErrNotStruct      = errors.New("not struct")
	ErrConflictPath   = errors.New("conflict path")
	ErrNotMultipleOf2 = errors.New("not multiple of 2")
)

const (
	POST MethodType = iota
	GET
	PATCH
	HEAD
	PUT
	DELETE
	OPTION
	CONNECT
	TRACE
)

const (
	ObjectTypeRequest ObjectType = iota
	ObjectTypeReply
	ObjectTypeObject
)

const (
	TmplTypeStruct TmplType = iota
	TmplTypeInterface
	TmplTypeEq
)
