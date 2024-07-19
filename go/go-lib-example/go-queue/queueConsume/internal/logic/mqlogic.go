package logic

import (
	"context"
	"fmt"
	"queueConsume/internal/config"
	"queueConsume/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

type MqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MqLogic {
	return &MqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MqLogic) Consume(_, val string) error {
	fmt.Println("consume msg", val)
	return nil
}

func Consumers(c config.Config, ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(c.KqConf, kq.WithHandle(func(k, v string) error {
			fmt.Printf("=>xxxx %s\n", v)
			return nil
		})),
		//kq.MustNewQueue(c.KqConf, NewMqLogic(ctx, svcCtx)),
		//kq.MustNewQueue(c.KqConf, NewMqLogic(ctx, svcCtx)),
	}
}
