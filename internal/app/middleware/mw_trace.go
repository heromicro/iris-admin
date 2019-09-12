package middleware

import (

	"github.com/wanhello/iris-admin/internal/app/irisplus"
	"github.com/wanhello/iris-admin/pkg/util"
	"github.com/kataras/iris"
)

// TraceMiddleware 跟踪ID中间件
func TraceMiddleware(skipper ...SkipperFunc) iris.HandlerFunc {
	return func(c *iris.Context) {
		if len(skipper) > 0 && skipper[0](c) {
			c.Next()
			return
		}

		// 优先从请求头中获取请求ID，如果没有则使用UUID
		traceID := c.GetHeader("X-Request-Id")
		if traceID == "" {
			traceID = util.MustUUID()
		}
		c.Set(irisplus.TraceIDKey, traceID)
		c.Next()
	}
}
