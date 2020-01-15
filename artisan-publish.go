package artisan

import (
	"fmt"
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
	"os"
)

type PublishedServices struct {
	svc         []ServiceDescription
	packageName string
	//wildObjTemplates []ObjTmpl
	wildSvc ServiceDescription
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
	for i := range c.svc {
		err = publish(c.packageName, c.svc[i])
		if err != nil {
			return
		}
	}
	err = publish(c.packageName, c.wildSvc)
	return
}

func publish(packageName string, svc ServiceDescription) (err error) {
	ctx := newTmplContext(svc)
	ctx.AppendPackage("github.com/Myriad-Dreamin/minimum-lib/controller")

	objs, funcs := svc.GenerateObjects(DefaultFunctionTmplFactories, ctx)

	//fmt.Println(packages)
	sugar.WithWriteFile(func(f *os.File) {
		_, err = fmt.Fprintf(f, `
package %s

import (
%s
)

var _ controller.MContext

%s`, packageName, depList(ctx.GetPackages()), svcIface(svc))
		if err != nil {
			return
		}

		for _, obj := range objs {
			_, err = f.WriteString(obj.String())
			if err != nil {
				return
			}
			_, err = f.WriteString("\n")
			if err != nil {
				return
			}
		}

		for _, v := range funcs {
			_, err = f.WriteString(v.String())
			if err != nil {
				return
			}
			_, err = f.WriteString("\n")
			if err != nil {
				return
			}
		}
	}, svc.GetFilePath())
	return
}

func svcIface(svc ServiceDescription) string {
	if len(svc.GetName()) == 0 {
		return ""
	}
	return fmt.Sprintf(`
type %s interface {
%s
}`, svc.GetName(), svcMethods(svc))
}

func svcMethods(svc ServiceDescription) (res string) {
	res = fmt.Sprintf("    %sSignatureXXX() interface{}\n", svc.GetName())
	for _, cat := range svc.GetCategories() {
		for _, method := range cat.GetMethods() {
			res += "    " + method.GetName() + "(c controller.MContext)\n"
		}
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
