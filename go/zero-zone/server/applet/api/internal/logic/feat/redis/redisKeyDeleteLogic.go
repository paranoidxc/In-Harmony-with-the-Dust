package redis

import (
	"context"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RedisKeyDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRedisKeyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedisKeyDeleteLogic {
	return &RedisKeyDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RedisKeyDeleteLogic) RedisKeyDelete(req *types.RedisKeyDeleteReq) error {
	_, err := l.svcCtx.Redis.Del(req.Key)
	return err
}
