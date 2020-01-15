package artisan

import (
	"math/big"
	"reflect"
)

type parameterDescription struct {
	embedObjects []ObjectDescription
	pType        Type
	name         string
	field        Field
	source       *source
	tags         map[string]string

	hash []byte
}

func (p *parameterDescription) Hash() []byte {
	if p.hash == nil {
		p.genHash()
	}
	return p.hash
}

func (p *parameterDescription) GetSource() *source {
	return p.source
}

func (p *parameterDescription) GetDTOName() string {
	return p.name
}

func (p *parameterDescription) GetField() Field {
	return p.field
}

func (p *parameterDescription) GetTag() TagI {
	return p.tags
}

func (p *parameterDescription) GetEmbedObjects() []ObjectDescription {
	return p.embedObjects
}

func (p *parameterDescription) GetType() Type {
	return p.pType
}

func (p *parameterDescription) genHash() {
	p.hash = big.NewInt(int64(reflect.ValueOf(p).Elem().UnsafeAddr())).Bytes()
}
