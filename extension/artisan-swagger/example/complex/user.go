package complex_example

import (
	"github.com/Myriad-Dreamin/artisan/artisan-core"
	"github.com/Myriad-Dreamin/artisan/extension/artisan-swagger/example/complex/model"
)

type UserCategories struct {
	artisan_core.VirtualService
	List           artisan_core.Category
	Login          artisan_core.Category
	Register       artisan_core.Category
	ChangePassword artisan_core.Category
	Inspect        artisan_core.Category
	IdGroup        artisan_core.Category
}

func DescribeUserService(base string) artisan_core.ProposingService {
	meta := Meta{
		RouterMeta: artisan_core.RouterMeta{RuntimeRouterMeta: "user"},
	}

	var userModel = new(model.User)
	var vUserModel model.User
	svc := &UserCategories{
		List: artisan_core.Ink().
			Path("user-list").
			Method(artisan_core.POST, "ListUsers",
				artisan_core.QT("ListUsersRequest", model.Filter{}),
				artisan_core.Reply(
					codeField,
					artisan_core.ArrayParam(artisan_core.Param("users", artisan_core.Object(
						"ListUserReply",
						artisan_core.SPsC(&vUserModel.NickName, &vUserModel.LastLogin),
					))),
				),
			),
		Login: artisan_core.Ink().
			Path("login").
			Method(artisan_core.POST, "Login",
				artisan_core.Request(
					artisan_core.SPsC(&userModel.ID, &userModel.NickName, &userModel.Phone),
					artisan_core.Param("password", artisan_core.String, required),
				),
				artisan_core.Reply(
					codeField,
					artisan_core.SPsC(&userModel.ID, &userModel.Phone, &userModel.NickName, &userModel.Name),
					artisan_core.Param("identity", artisan_core.Strings),
					artisan_core.Param("token", artisan_core.String),
					artisan_core.Param("refresh_token", artisan_core.String),
				),
			),
		Register: artisan_core.Ink().
			Path("register").
			Method(artisan_core.POST, "Register",
				artisan_core.Request(
					artisan_core.SPs(artisan_core.C(&userModel.Name, &userModel.NickName, &userModel.Phone), required),
					artisan_core.Param("password", artisan_core.String, required),
				),
				artisan_core.Reply(
					codeField,
					artisan_core.Param("id", &userModel.ID)),
			),
		ChangePassword: artisan_core.Ink().
			Path("user/:id/password").
			Method(artisan_core.PUT, "ChangePassword",
				artisan_core.Request(
					artisan_core.Param("old_password", artisan_core.String, required),
					artisan_core.Param("new_password", artisan_core.String, required),
				),
			),
		Inspect: artisan_core.Ink().Path("user/:id/inspect").
			Method(artisan_core.GET, "InspectUser",
				artisan_core.Reply(
					codeField,
					artisan_core.Param("user", &userModel),
				),
			),
		IdGroup: artisan_core.Ink().
			Path("user/:id").
			Method(artisan_core.GET, "GetUser",
				artisan_core.Reply(
					codeField,
					artisan_core.SPsC(&userModel.NickName, &userModel.LastLogin),
				)).
			Method(artisan_core.PUT, "PutUser",
				artisan_core.Request(
					artisan_core.Param("phone", &userModel.Phone),
				)).
			Method(artisan_core.DELETE, "Delete"),
	}
	svc.Name("UserService").Base(base).Meta(meta).UseModel(
		artisan_core.Model(artisan_core.Name("user"), &userModel),
		artisan_core.Model(artisan_core.Name("vUser"), &vUserModel))
	return svc
}
