package redis

import (
	"context"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RedisKeyDeletesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRedisKeyDeletesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedisKeyDeletesLogic {
	return &RedisKeyDeletesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RedisKeyDeletesLogic) RedisKeyDeletes(req *types.RedisKeyDeletesReq) error {
	_, err := l.svcCtx.Redis.Del(req.Key...)
	return err
}
