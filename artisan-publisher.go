package artisan

import (
	"reflect"
)

type PublishingServices struct {
	rawSvc      []ProposingService
	packageName string
}

func (c *PublishingServices) AppendService(rawSvc ...ProposingService) *PublishingServices {
	c.rawSvc = append(c.rawSvc, rawSvc...)
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
	for _, svc := range c.rawSvc {

		// compile models
		ctx := &Context{
			svc: svc,
		}
		ctx.makeSources()

		// get name and file path of service
		desc := &serviceDescription{
			name:     svc.GetName(),
			base:     svc.GetBase(),
			filePath: svc.GetFilePath(),
		}

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

		d.svc = append(d.svc, desc)
	}
	return d
}
