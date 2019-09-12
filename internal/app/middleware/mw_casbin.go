package middleware

import (
	"github.com/wanhello/iris-admin/internal/app/config"
	"github.com/wanhello/iris-admin/internal/app/errors"
	"github.com/wanhello/iris-admin/internal/app/irisplus"

	"github.com/casbin/casbin"
	"github.com/kataras/iris"

)

// CasbinMiddleware casbin中间件
func CasbinMiddleware(enforcer *casbin.Enforcer, skipper ...SkipperFunc) iris.HandlerFunc {
	cfg := config.GetGlobalConfig()
	return func(c *iris.Context) {
		if !cfg.EnableCasbin || len(skipper) > 0 && skipper[0](c) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		m := c.Request.Method
		if b, err := enforcer.EnforceSafe(irisplus.GetUserID(c), p, m); err != nil {
			irisplus.ResError(c, errors.WithStack(err))
			return
		} else if !b {
			irisplus.ResError(c, errors.ErrNoResourcePerm)
			return
		}
		c.Next()
	}
}
