package middleware

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/kataras/iris"
)

// WWWMiddleware 静态站点中间件
func WWWMiddleware(root string, skipper ...SkipperFunc) iris.HandlerFunc {
	return func(c iris.Context) {
		if len(skipper) > 0 && skipper[0](c) {
			c.Next()
			return
		}

		p := c.Request().URL.Path
		fpath := filepath.Join(root, filepath.FromSlash(p))
		_, err := os.Stat(fpath)
		if err != nil && os.IsNotExist(err) {
			fpath = filepath.Join(root, "index.html")
		}

		http.ServeFile(c.ResponseWriter(), c.Request(), fpath)
		c.Abort()
	}
}
