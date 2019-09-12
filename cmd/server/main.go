package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/wanhello/iris-admin/internal/app"
	"github.com/wanhello/iris-admin/pkg/logger"
	"github.com/wanhello/iris-admin/pkg/util"
)


// VERSION 版本号，
// 可以通过编译的方式指定版本号：go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "1.0.0"


var (
	configFile string
	modelFile  string
	wwwDir     string
	swaggerDir string
)


func init() {
	flag.StringVar(&configFile, "c", "", "config_file(.json,.yaml,.toml)")
	flag.StringVar(&modelFile, "m", "", "Casbin access model config(.conf)")
	flag.StringVar(&wwwDir, "www", "", "static directory")
	flag.StringVar(&swaggerDir, "swagger", "", "swagger directory")
}


func main() {
	flag.Parse()

	if configFile == "" {
		panic("请使用-c指定配置文件")
	}

	var state int32 = 1
	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx := logger.NewTraceIDContext(context.Background(), util.MustUUID())
	span := logger.StartSpanWithCall(ctx)

	call := app.Init(ctx,
		app.SetConfigFile(configFile),
		app.SetModelFile(modelFile),
		app.SetWWWDir(wwwDir),
		app.SetSwaggerDir(swaggerDir),
		app.SetVersion(VERSION))

	select {
	case sig := <-sc:
		atomic.StoreInt32(&state, 0)
		span().Printf("获取到退出信号[%s]", sig.String())
	}

	if call != nil {
		call()
	}
	span().Printf("服务退出")

	os.Exit(int(atomic.LoadInt32(&state)))
}







