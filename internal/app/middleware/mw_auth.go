package middleware

import (

	"github.com/wanhello/iris-admin/internal/app/config"
	"github.com/wanhello/iris-admin/internal/app/errors"

	"github.com/wanhello/iris-admin/internal/app/irisplus"

	"github.com/wanhello/iris-admin/pkg/auth"

	"github.com/kataras/iris"
)



// UserAuthMiddleware 用户授权中间件
func UserAuthMiddleware(a auth.Auther, skipper ...SkipperFunc) iris.HandlerFunc {
	return func(c *iris.Context) {
		var userID string
		if t := irisplus.GetToken(c); t != "" {
			id, err := a.ParseUserID(t)
			if err != nil {
				if err == auth.ErrInvalidToken {
					irisplus.ResError(c, errors.ErrNoPerm)
					return
				}
				irisplus.ResError(c, errors.WithStack(err))
				return
			}
			userID = id
		}

		if userID != "" {
			c.Set(irisplus.UserIDKey, userID)
		}

		if len(skipper) > 0 && skipper[0](c) {
			c.Next()
			return
		}

		if userID == "" {
			if config.GetGlobalConfig().RunMode == "debug" {
				c.Set(irisplus.UserIDKey, config.GetGlobalConfig().Root.UserName)
				c.Next()
				return
			}
			irisplus.ResError(c, errors.ErrNoPerm)
		}
	}
}


