package logic

import (
	"context"

	"queueapi/internal/svc"
	"queueapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueueapiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueueapiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueueapiLogic {
	return &QueueapiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueueapiLogic) Queueapi(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
