package main

import (
	"flag"
	"fmt"
	"net/http"
	"zero-zone/applet/api/internal/cronjob"
	"zero-zone/pkg/errorx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"zero-zone/applet/api/internal/config"
	"zero-zone/applet/api/internal/handler"
	"zero-zone/applet/api/internal/svc"
)

var configFile = flag.String("f", "etc/core-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.OldRegisterHandlers(server, ctx)
	handler.RegisterHandlers(server, ctx)

	cronjob.InitTimer(server, ctx)

	// 自定义错误
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case *errorx.CodeError:
			return http.StatusOK, e.Data()
		default:
			return http.StatusInternalServerError, nil
		}
	})

	if c.Mode == "dev" {
		//logx.DisableStat()
	}

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
