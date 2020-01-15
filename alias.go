package artisan

func Q(desc ...interface{}) RequestObject {
	return Request(desc...)
}

func T(name string, base interface{}) *transferClass {
	return Transfer(name, base)
}

func QT(name string, base interface{}) RequestObject {
	return Q(T(name, base))
}

func SP(descriptions ...interface{}) Parameter {
	return SnakeParam(descriptions...)
}

func P(name string, descriptions ...interface{}) Parameter {
	return Param(name, descriptions...)
}

func C(pack ...interface{}) []interface{} {
	return pack
}

func SPs(mTypes []interface{}, descriptions ...interface{}) (p []Parameter) {
	for i := range mTypes {
		p = append(p, SP(append(descriptions, mTypes[i])...))
	}
	return
}

func SPsC(mTypes ...interface{}) []Parameter {
	return SPs(mTypes)
}

func Ps(mNameTypes []interface{}, descriptions ...interface{}) (p []Parameter) {
	if (len(mNameTypes) & 1) != 0 {
		panic(ErrNotMultipleOf2)
	}
	for i := 0; i < len(mNameTypes); i += 2 {
		name, t := mNameTypes[i], mNameTypes[i+1]
		p = append(p, P(name.(string), append(descriptions, t)...))
	}
	return
}
