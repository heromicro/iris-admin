package middleware

import (
	"time"

	"github.com/wanhello/iris-admin/internal/app/config"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"

)

// CORSMiddleware 跨域请求中间件
func CORSMiddleware() iris.HandlerFunc {
	cfg := config.GetGlobalConfig().CORS
	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           time.Second * time.Duration(cfg.MaxAge),
	})
}

