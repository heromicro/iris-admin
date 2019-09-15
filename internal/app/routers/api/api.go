package api

import (
	"github.com/casbin/casbin"
	"github.com/kataras/iris"
	"github.com/wanhello/iris-admin/internal/app/middleware"
	"github.com/wanhello/iris-admin/internal/app/routers/api/ctl"
	"github.com/wanhello/iris-admin/pkg/auth"

	"go.uber.org/dig"

)


// RegisterRouter 注册/api路由
func RegisterRouter(app *iris.Application, container *dig.Container) error {
	err := ctl.Inject(container)
	if err != nil {
		return err
	}

	return container.Invoke(func(
		a auth.Auther,
		e *casbin.Enforcer,
		cDemo *ctl.Demo,
		cLogin *ctl.Login,
		cMenu *ctl.Menu,
		cRole *ctl.Role,
		cUser *ctl.User,
	) error {

		g := app.Party("/api")

		// 用户身份授权
		g.Use(middleware.UserAuthMiddleware(
			a,
			middleware.AllowMethodAndPathPrefixSkipper(
				middleware.JoinRouter("GET", "/api/v1/pub/login"),
				middleware.JoinRouter("POST", "/api/v1/pub/login"),
			),
		))

		// casbin权限校验中间件
		g.Use(middleware.CasbinMiddleware(e,
			middleware.AllowMethodAndPathPrefixSkipper(
				middleware.JoinRouter("GET", "/api/v1/pub"),
				middleware.JoinRouter("POST", "/api/v1/pub"),
			),
		))

		// 请求频率限制中间件
		g.Use(middleware.RateLimiterMiddleware())

		v1 := g.Party("/v1")
		{
			pub := v1.Party("/pub")
			{
				// 注册/api/v1/pub/login
				pub.Get("/login/captchaid", cLogin.GetCaptcha)
				pub.Get("/login/captcha", cLogin.ResCaptcha)
				pub.Post("/login", cLogin.Login)
				pub.Post("/login/exit", cLogin.Logout)

				// 注册/api/v1/pub/refresh_token
				pub.Post("/refresh_token", cLogin.RefreshToken)

				// 注册/api/v1/pub/current
				pub.Put("/current/password", cLogin.UpdatePassword)
				pub.Get("/current/user", cLogin.GetUserInfo)
				pub.Get("/current/menutree", cLogin.QueryUserMenuTree)
			}

			// 注册/api/v1/demos
			v1.Get("/demos", cDemo.Query)
			v1.Get("/demos/:id", cDemo.Get)
			v1.Post("/demos", cDemo.Create)
			v1.Put("/demos/:id", cDemo.Update)
			v1.Delete("/demos/:id", cDemo.Delete)
			v1.Patch("/demos/:id/enable", cDemo.Enable)
			v1.Patch("/demos/:id/disable", cDemo.Disable)

			// 注册/api/v1/menus
			v1.Get("/menus", cMenu.Query)
			v1.Get("/menus/:id", cMenu.Get)
			v1.Post("/menus", cMenu.Create)
			v1.Put("/menus/:id", cMenu.Update)
			v1.Delete("/menus/:id", cMenu.Delete)

			// 注册/api/v1/roles
			v1.Get("/roles", cRole.Query)
			v1.Get("/roles/:id", cRole.Get)
			v1.Post("/roles", cRole.Create)
			v1.Put("/roles/:id", cRole.Update)
			v1.Delete("/roles/:id", cRole.Delete)

			// 注册/api/v1/users
			v1.Get("/users", cUser.Query)
			v1.Get("/users/:id", cUser.Get)
			v1.Post("/users", cUser.Create)
			v1.Put("/users/:id", cUser.Update)
			v1.Delete("/users/:id", cUser.Delete)
			v1.Patch("/users/:id/enable", cUser.Enable)
			v1.Patch("/users/:id/disable", cUser.Disable)
		}

		return nil
	})
}
