package artisan

import (
	"reflect"
)

type wildService struct {
	VirtualService
	WildCate Category
}

func newWildService() *wildService {
	return &wildService{
		WildCate: newCategory(),
	}
}

type PublishingServices struct {
	rawSvc      []ProposingService
	packageName string

	mixWildModel bool
	wildSvc      *wildService
}

func (c *PublishingServices) AppendService(rawSvc ...ProposingService) *PublishingServices {
	c.rawSvc = append(c.rawSvc, rawSvc...)
	return c
}

func (c *PublishingServices) MixWildModel() *PublishingServices {
	c.mixWildModel = true
	return c
}

func (c *PublishingServices) UseModel(models ...*model) *PublishingServices {
	c.wildSvc.UseModel(models...)
	return c
}

func (c *PublishingServices) WildPath(filePath string) *PublishingServices {
	c.wildSvc.ToFile(filePath)
	return c
}

func (c *PublishingServices) Object(descriptions ...interface{}) *PublishingServices {
	c.wildSvc.WildCate = c.wildSvc.WildCate.HelpWrapObjectXXX(1, descriptions...)
	return c
}

func (c *PublishingServices) AppendObject(objs ...SerializeObject) *PublishingServices {
	c.wildSvc.WildCate = c.wildSvc.WildCate.AppendObject(objs...)
	return c
}

func (c *PublishingServices) SetPackageName(packageName string) *PublishingServices {
	c.packageName = packageName
	return c
}

func (c *PublishingServices) GetPackageName() string {
	return c.packageName
}

func (c *PublishingServices) Publish() error {
	return c.Final().Publish()
}

func (c *PublishingServices) Final() (d *PublishedServices) {
	d = new(PublishedServices)
	d.packageName = c.packageName
	d.wildSvc = makeServiceDescription(c.wildSvc)
	for _, svc := range c.rawSvc {
		if c.mixWildModel {
			svc.UseModel(c.wildSvc.GetModels()...)
		}
		d.svc = append(d.svc, makeServiceDescription(svc))
	}
	return d
}

func makeServiceDescription(svc ProposingService) *serviceDescription {
	// compile models
	ctx := &Context{
		rawSvc: svc,
	}
	// get name and file path of service
	desc := &serviceDescription{
		name:     svc.GetName(),
		base:     svc.GetBase(),
		filePath: svc.GetFilePath(),
	}
	ctx.svc = desc

	ctx.makeSources()

	// build category methods
	value, svcType := getElements(svc)
	if svcType.Kind() != reflect.Struct {
		panic(ErrNotStruct)
	}
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if cate, ok := field.Interface().(Category); ok && cate != nil {
			desc.categories = append(desc.categories,
				cate.CreateCategoryDescription(ctx))
		}
	}

	//// get packages
	//desc.packages = ctx.packages
	return desc
}