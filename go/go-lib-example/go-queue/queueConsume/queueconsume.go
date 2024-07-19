package main

import (
	"context"
	"flag"
	"fmt"

	"queueConsume/internal/config"
	"queueConsume/internal/handler"
	"queueConsume/internal/logic"
	"queueConsume/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	zeroservice "github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/queueconsume-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	svcCtx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, svcCtx)

	/*
		q := kq.MustNewQueue(c.KqConf, kq.WithHandle(func(k, v string) error {
			fmt.Printf("=> %s\n", v)
			return nil
		}))
		defer q.Stop()
		q.Start()
	*/

	ctx := context.Background()

	serviceGroup := zeroservice.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range logic.Consumers(c, ctx, svcCtx) {
		fmt.Println("mq", mq)
		serviceGroup.Add(mq)
	}
	serviceGroup.Start()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
