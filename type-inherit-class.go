package artisan

type inheritClass struct {
	uuid  UUID
	name  string
	bases []interface{}
	dp    string
}

func (i inheritClass) DefiningPosition() string {
	return i.dp
}

func (i inheritClass) GenObjectTmpl() ObjTmpl {
	panic("implement me")
}

func (i inheritClass) GetType() Type {
	panic("implement me")
}

func (i inheritClass) GetUUID() UUID {
	return i.uuid
}

func (i inheritClass) String() string {
	panic("implement me")
	//return fmt.Sprintf("type %s = %s", i.name, reflect.TypeOf(i.base))
}

func (i inheritClass) GetTypeString() string {
	return i.name
}

func (i inheritClass) GetEmbedObject() []ObjectDescription {
	return nil
}

func (i inheritClass) CreateObjectDescription(ctx *Context) ObjectDescription {
	return i
}
