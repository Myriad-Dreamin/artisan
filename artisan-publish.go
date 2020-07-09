package artisan

type PublishedServices struct {
	svcMap      map[ProposingService]ServiceDescription
	packageName string
	//wildObjTemplates []ObjTmpl
	wildSvc ServiceDescription

	Opts *PublishOptions
}

func (c *PublishedServices) PublishInterface(svc ServiceDescription) error {
	return svc.PublishInterface(c.packageName, c.Opts)
}

func (c *PublishedServices) PublishObjects(svc ServiceDescription) error {
	return svc.PublishObjects(c.packageName, c.Opts)
}

func (c *PublishedServices) SetOptions(opts *PublishOptions) *PublishedServices {
	c.Opts = opts
	return c
}

func (c *PublishedServices) GetPackageName() string {
	return c.packageName
}

func (c *PublishedServices) GetService(raw ProposingService) ServiceDescription {
	return c.svcMap[raw]
}

func (c *PublishedServices) GetServices() map[ProposingService]ServiceDescription {
	return c.svcMap
}

func (c *PublishedServices) GetWildServices() ServiceDescription {
	return c.wildSvc
}

func (c *PublishedServices) Publish() error {
	if err := c.writeToFiles(); err != nil {
		return err
	}
	return nil
}

func (c *PublishedServices) writeToFiles() (err error) {
	if err = c.writeSVCsAndDTOs(); err != nil {
		return
	}
	return
}

func (c *PublishedServices) writeSVCsAndDTOs() (err error) {
	for _, svc := range c.svcMap {
		err = svc.PublishAll(c.packageName, c.Opts)
		if err != nil {
			return
		}
	}
	if len(c.wildSvc.GetFilePath()) != 0 {
		err = c.wildSvc.PublishAll(c.packageName, c.Opts)
	}
	return
}

func depList(pkgSet map[string]bool) (res string) {
	for k := range pkgSet {
		if len(k) > 0 {
			res += `    "` + k + `"
`
		}
	}
	return
}
