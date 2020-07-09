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
	base        string

	mixWildModel bool
	wildSvc      *wildService

	Opts *PublishOptions
}

func (c *PublishingServices) GetRawProtocols() []ProposingService {
	return append(c.rawSvc, c.wildSvc)
}

func (c *PublishingServices) GetBase() string {
	return c.base
}

func (c *PublishingServices) GetPackageName() string {
	return c.packageName
}

func (c *PublishingServices) IsMixWildModel() bool {
	return c.mixWildModel
}

func (c *PublishingServices) GetWildModels() []*model {
	return c.wildSvc.GetModels()
}

func (c *PublishingServices) GetWildName() string {
	return c.wildSvc.GetName()
}

func (c *PublishingServices) GetWildBase() string {
	return c.wildSvc.GetBase()
}

func (c *PublishingServices) GetWildFilePath() string {
	return c.wildSvc.GetFilePath()
}

func (c *PublishingServices) SetOptions(opts *PublishOptions) *PublishingServices {
	c.Opts = opts
	return c
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

func (c *PublishingServices) Base(urlPath string) *PublishingServices {
	c.base = urlPath
	return c
}

func (c *PublishingServices) WildToFile(filePath string) *PublishingServices {
	c.wildSvc.ToFile(filePath)
	return c
}

func (c *PublishingServices) WildBase(urlPath string) *PublishingServices {
	c.wildSvc.Base(urlPath)
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

func (c *PublishingServices) Publish() error {
	return c.Final().Publish()
}

func (c *PublishingServices) Final() (d *PublishedServices) {
	d = new(PublishedServices)
	d.svcMap = make(map[ProposingService]ServiceDescription)
	d.packageName = c.packageName
	d.wildSvc = makeServiceDescription(c.wildSvc)
	d.Opts = c.Opts
	for _, svc := range c.rawSvc {
		if c.mixWildModel {
			svc.UseModel(c.wildSvc.GetModels()...)
		}
		d.svcMap[svc] = makeServiceDescription(svc)
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
		meta:     svc.GetMeta(),
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
			cd := cate.CreateCategoryDescription(ctx)
			if len(cd.GetName()) == 0 {
				cd.SetName(svcType.Field(i).Name)
			}
			desc.categories = append(desc.categories, cd)
		}
	}

	//// get packages
	//desc.packages = ctx.packages
	return desc
}
