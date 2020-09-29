package artisan_core

import (
	"fmt"
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
	"os"
	"strings"
)

type PublishOptions struct {
	InterfaceStyle string
	RouterFilePath string
}

type PublishedServices struct {
	HumanInfo   interface{}
	SvcMap      map[ProposingService]ServiceDescription
	PackageName string
	//wildObjTemplates []ObjTmpl
	WildSvc ServiceDescription

	Opts *PublishOptions
}

func (c *PublishedServices) PublishInterface(svc ServiceDescription) error {
	return svc.PublishInterface(c.PackageName, c.Opts)
}

func (c *PublishedServices) PublishObjects(svc ServiceDescription) error {
	return svc.PublishObjects(c.PackageName, c.Opts)
}

func (c *PublishedServices) SetOptions(opts *PublishOptions) *PublishedServices {
	c.Opts = opts
	return c
}

func (c *PublishedServices) GetPackageName() string {
	return c.PackageName
}

func (c *PublishedServices) GetService(raw ProposingService) ServiceDescription {
	return c.SvcMap[raw]
}

func (c *PublishedServices) GetServices() map[ProposingService]ServiceDescription {
	return c.SvcMap
}

func (c *PublishedServices) GetWildServices() ServiceDescription {
	return c.WildSvc
}

func (c *PublishedServices) Publish() error {
	if err := c.writeToFiles(); err != nil {
		return err
	}
	return nil
}

func (c *PublishedServices) PublishRouter(filePath string) (err error) {

	defs, invokes, fs := c.getSVCRouterFragments()

	sugar.WithWriteFile(func(f *os.File) {
		_, err = fmt.Fprintf(f, `
package %s

import (
    "github.com/Myriad-Dreamin/minimum-lib/controller"
%s
)

type Router = controller.Router
type Middleware = controller.Middleware
type LeafRouter = controller.LeafRouter
type HandlerFunc = controller.HandlerFunc

type H interface {
    GetRouter() *Router
    GetAuthRouter() *Router
    GetAuth() *Middleware
}

type BaseH struct {
    *Router
    AuthRouter *Router
    Auth       *Middleware
}

func (r *BaseH) GetRouter() *Router {
    return r.Router
}

func (r *BaseH) GetAuthRouter() *Router {
    return r.AuthRouter
}

func (r *BaseH) GetAuth() *Middleware {
    return r.Auth
}

type GenerateRouterTraits interface {
    GetJWTMiddleware() HandlerFunc
    GetAuthMiddleware() *Middleware
    AfterBuild(r *RootRouter)
    ApplyAuth(r *RootRouter)
    ApplyAuthOnMethod(r *LeafRouter, authMeta string) *LeafRouter

    ApplyRouteMeta(m *Middleware, routeMeta string) *Middleware
    GetServiceInstance(svcName string) interface{}
}

type RootRouter struct {
    H
    Root *Router
    %s
    Ping     *LeafRouter
    //Images   *LeafRouter
    //Musics   *LeafRouter
    //Articles *LeafRouter
}

// @title Ping
// @description result
func PingFunc(c controller.MContext) {
    c.JSON(200, map[string]interface{}{
        "message": "pong",
    })
}

%s
%s`, c.PackageName, "", defs, c.genRootRouterFunc(invokes), fs)
		if err != nil {
			return
		}
	}, filePath)
	return
}

func (c *PublishedServices) writeToFiles() (err error) {
	if err = c.writeSVCsAndDTOs(); err != nil {
		return
	}
	return
}

func (c *PublishedServices) writeSVCsAndDTOs() (err error) {
	for _, svc := range c.SvcMap {
		err = svc.PublishAll(c.PackageName, c.Opts)
		if err != nil {
			return
		}
	}
	if len(c.WildSvc.GetFilePath()) != 0 {
		err = c.WildSvc.PublishAll(c.PackageName, c.Opts)
	}
	return
}

func (c *PublishedServices) genRootRouterFunc(invokes string) string {
	return fmt.Sprintf(`
func NewRootRouter(traits GenerateRouterTraits) (r *RootRouter) {
    rr := controller.NewRouterGroup()
    apiRouterV1 := rr.Group("/v1")
    authRouterV1 := apiRouterV1.Group("", traits.GetJWTMiddleware())

    r = &RootRouter{
        Root: rr,
        H: &BaseH{
            Router:     apiRouterV1,
            AuthRouter: authRouterV1,
            Auth:       traits.GetAuthMiddleware(),
        },
    }

    r.Ping = r.Root.GET("/ping", PingFunc)

    %s

    traits.AfterBuild(r)
    traits.ApplyAuth(r)
    return
}`, invokes)
}

func (c *PublishedServices) getSVCRouterFragments() (string, string, string) {
	if c.Opts == nil {
		c.Opts = &PublishOptions{}
	}

	var defs, invokes, fnDefs []string
	for _, svc := range c.SvcMap {
		ctx := newTmplContext(svc)
		sn, fn := svc.GenerateRouterFunc(ctx, c.Opts.InterfaceStyle)

		defs = append(defs, fmt.Sprintf("%s *%s", svc.GetName(), sn))
		invokes = append(invokes, fmt.Sprintf("r.%s = New%s(traits, r.H)", svc.GetName(), sn))

		fnDefs = append(fnDefs, fn)
	}

	return strings.Join(defs, "\n    "), strings.Join(invokes, "\n    "),
		strings.Join(fnDefs, "\n")
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

func routeGen(ctx TmplCtx, interfaceStyle string, node GenTreeNode) (
	string, string, string) {

	var defs, invokes, fnDefs []string
	for _, svc := range node.GetCategories() {
		sn, fn := svc.GenerateRouterFunc(ctx, interfaceStyle)

		defs = append(defs, fmt.Sprintf("%s *%s", svc.GetName(), sn))
		invokes = append(invokes, fmt.Sprintf("r.%s = New%s(traits, r.H)", svc.GetName(), sn))

		fnDefs = append(fnDefs, fn)
	}

	return strings.Join(defs, "\n    "), strings.Join(invokes, "\n    "),
		strings.Join(fnDefs, "\n")
}

func GenTreeNodeRouteGen(ctx TmplCtx, svc GenTreeNode, interfaceStyle string) (string, string) {
	switch interfaceStyle {
	case InterfaceStyleMinimum:
		fallthrough
	default:
		if len(svc.GetName()) == 0 {
			return "", ""
		}
		var copyRouteInvoke string
		if len(svc.GetBase()) != 0 {
			copyRouteInvoke = fmt.Sprintf(`Group("%s")`, svc.GetBase())
		} else {
			copyRouteInvoke = fmt.Sprintf(`Extend("%s")`, svc.GetName())
		}

		rsn := ctx.Get(VarContextRouteStructName)

		var sns string
		if x, ok := rsn.(string); ok {
			sns = x
		}
		ctx.Set(VarContextRouteStructName, sns+svc.GetName())

		var def, invoke, f = routeGen(ctx, interfaceStyle, svc)
		var methodDefs, methods []string

		var rm string
		if irm, ok := svc.GetMeta().(IRouterMeta); ok {
			rm = irm.GetRuntimeRouterMeta()
		}

		if g, ok := svc.(GenTreeNodeWithMethods); ok {
			for _, m := range g.GetMethods() {
				var sub = "GetRouter()"
				if len(m.GetAuthMeta()) != 0 {
					sub = "GetAuthRouter()"
				}

				methodDefs = append(methodDefs, fmt.Sprintf(
					"%s *LeafRouter", m.GetName()))
				methods = append(methods,
					fmt.Sprintf(`r.%s = r.%s.%s("", traits.GetServiceInstance("%s").(%s).%s)`,
						m.GetName(),
						sub,
						MethodTypeMapping[m.GetMethodType()],
						ctx.Get(VarContextServiceStructName).(string),
						ctx.Get(VarContextServiceStructName).(string),
						m.GetName()),
				)
				if len(m.GetAuthMeta()) != 0 {
					methods = append(methods,
						fmt.Sprintf(`r.%s = traits.ApplyAuthOnMethod(r.%s, "%s")`,
							m.GetName(),
							m.GetName(),
							m.GetAuthMeta(),
						),
					)
				}
			}
		}

		ctx.Set(VarContextRouteStructName, rsn)
		sn := fmt.Sprintf("%sRouter", sns+svc.GetName())

		return sn, fmt.Sprintf(`
type %s struct {
    H
    %s

    %s
}


func New%s(traits GenerateRouterTraits, h H) (r *%s) {
    r = &%s{
        H: &BaseH{
            Router:     h.GetRouter().%s,
            AuthRouter: h.GetAuthRouter().%s,
            Auth:       traits.ApplyRouteMeta(h.GetAuth(), "%s"),
        },
    }

    %s

    %s

    return
}

%s`, sn, def, strings.Join(methodDefs, "\n    "), sn, sn, sn,
			copyRouteInvoke, copyRouteInvoke, rm, invoke, strings.Join(methods, "\n    "), f)
	}
}
