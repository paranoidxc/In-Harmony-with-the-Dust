package logic

import (
	"context"

	"queueConsume/internal/svc"
	"queueConsume/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueueConsumeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueueConsumeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueueConsumeLogic {
	return &QueueConsumeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueueConsumeLogic) QueueConsume(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
