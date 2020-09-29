package artisan_core

import (
	"reflect"
)

type WildService struct {
	VirtualService
	WildCate Category
}

func newWildService() *WildService {
	return &WildService{
		WildCate: newCategory(),
	}
}

type PublishingServices struct {
	humanInfo interface{}

	packageName string
	base        string

	mixWildModel bool

	RawSvc  []ProposingService
	WildSvc *WildService
	Opts    *PublishOptions
}

func (c *PublishingServices) GetHumanInfo() interface{} {
	return c.humanInfo
}

func (c *PublishingServices) HumanInfo(humanInfo interface{}) *PublishingServices {
	c.humanInfo = humanInfo
	return c
}

func (c *PublishingServices) GetRawProtocols() []ProposingService {
	return append(c.RawSvc, c.WildSvc)
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
	return c.WildSvc.GetModels()
}

func (c *PublishingServices) GetWildName() string {
	return c.WildSvc.GetName()
}

func (c *PublishingServices) GetWildBase() string {
	return c.WildSvc.GetBase()
}

func (c *PublishingServices) GetWildFilePath() string {
	return c.WildSvc.GetFilePath()
}

func (c *PublishingServices) SetOptions(opts *PublishOptions) *PublishingServices {
	c.Opts = opts
	return c
}

func (c *PublishingServices) AppendService(rawSvc ...ProposingService) *PublishingServices {
	c.RawSvc = append(c.RawSvc, rawSvc...)
	return c
}

func (c *PublishingServices) MixWildModel() *PublishingServices {
	c.mixWildModel = true
	return c
}

func (c *PublishingServices) UseModel(models ...*model) *PublishingServices {
	c.WildSvc.UseModel(models...)
	return c
}

func (c *PublishingServices) Base(urlPath string) *PublishingServices {
	c.base = urlPath
	return c
}

func (c *PublishingServices) WildToFile(filePath string) *PublishingServices {
	c.WildSvc.ToFile(filePath)
	return c
}

func (c *PublishingServices) WildBase(urlPath string) *PublishingServices {
	c.WildSvc.Base(urlPath)
	return c
}

func (c *PublishingServices) Object(descriptions ...interface{}) *PublishingServices {
	c.WildSvc.WildCate = c.WildSvc.WildCate.HelpWrapObjectXXX(1, descriptions...)
	return c
}

func (c *PublishingServices) AppendObject(objs ...SerializeObject) *PublishingServices {
	c.WildSvc.WildCate = c.WildSvc.WildCate.AppendObject(objs...)
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
	d.HumanInfo = c.humanInfo
	d.PackageName = c.packageName
	d.Base = c.base
	d.Opts = c.Opts

	d.SvcMap = make(map[ProposingService]ServiceDescription)
	d.WildSvc = makeServiceDescription(c.WildSvc)
	for _, svc := range c.RawSvc {
		if c.mixWildModel {
			svc.UseModel(c.WildSvc.GetModels()...)
		}
		d.SvcMap[svc] = makeServiceDescription(svc)
	}
	return d
}

var dynCateType = reflect.TypeOf(new(Category)).Elem()

func makeServiceDescription(svc ProposingService) *serviceDescription {
	// compile models
	ctx := &Context{
		rawSvc: svc,
	}
	// get name and file path of service
	desc := &serviceDescription{
		humanInfo: svc.GetHumanInfo(),
		name:      svc.GetName(),
		base:      svc.GetBase(),
		meta:      svc.GetMeta(),
		models:    svc.GetModels(),
		filePath:  svc.GetFilePath(),
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
		if !field.Type().Implements(dynCateType) {
			continue
		}
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
