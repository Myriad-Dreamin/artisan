package artisan

import (
	"math/big"
	"reflect"
)

type parameterDescription struct {
	embedObjects       []ObjectDescription
	pType              Type
	name               string
	field              Field
	source             *source
	tags               map[string]string
	calculatedPackages PackageSet

	hash []byte
}

func (p *parameterDescription) Hash() []byte {
	if p.hash == nil {
		p.genHash()
	}
	return p.hash
}

func (p *parameterDescription) GetEmbedObjects() []ObjectDescription {
	return p.embedObjects
}

func (p *parameterDescription) GetType() Type {
	return p.pType
}

func (p *parameterDescription) GetDTOName() string {
	return p.name
}

func (p *parameterDescription) GetField() Field {
	return p.field
}

func (p *parameterDescription) GetSource() *source {
	return p.source
}

func (p *parameterDescription) GetTag() TagI {
	return p.tags
}

func (p *parameterDescription) GetPackages() PackageSet {
	var ps = make(PackageSet)
	for _, obj := range p.embedObjects {
		PackageSetInPlaceMerge(ps, obj.GetPackages())
	}
	return PackageSetInPlaceMerge(ps, p.calculatedPackages)
}

func (p *parameterDescription) genHash() {
	p.hash = big.NewInt(int64(reflect.ValueOf(p).Elem().UnsafeAddr())).Bytes()
}
