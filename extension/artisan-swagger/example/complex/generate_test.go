package complex_example

import (
	artisan_core "github.com/Myriad-Dreamin/artisan/artisan-core"
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
	"testing"
)

func TestGenerate(t *testing.T) {
	v1 := "v1"

	userCate := DescribeUserService("")
	objectCate := DescribeObjectService("")
	objectCate.ToFile("control/object.go")

	svc := artisan_core.NewService(
		userCate, objectCate).Base(v1).Final()

	sugar.HandlerError0(svc.PublishRouter("control/router.go"))

	userSvc := svc.GetService(userCate)
	delete(svc.GetServices(), userCate)

	sugar.HandlerError0(svc.PublishInterface(userSvc.SetFilePath("control/user-interface.go")))
	sugar.HandlerError0(svc.PublishObjects(userSvc.SetFilePath("control/user-objects.go")))

	sugar.HandlerError0(svc.Publish())
}
