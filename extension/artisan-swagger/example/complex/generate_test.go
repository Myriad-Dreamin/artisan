package complex_example

import (
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
	"testing"
)

func TestGenerate(t *testing.T) {
	svc := Generate()

	sugar.HandlerError0(svc.PublishRouter("control/router.go"))

	userSvc := svc.GetService(UserCate)
	delete(svc.GetServices(), UserCate)

	sugar.HandlerError0(svc.PublishInterface(userSvc.SetFilePath("control/user-interface.go")))
	sugar.HandlerError0(svc.PublishObjects(userSvc.SetFilePath("control/user-objects.go")))

	sugar.HandlerError0(svc.Publish())
}
