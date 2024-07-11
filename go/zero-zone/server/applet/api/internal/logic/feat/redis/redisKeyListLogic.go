package redis

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"sort"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type RedisKeyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRedisKeyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedisKeyListLogic {
	return &RedisKeyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RedisKeyListLogic) RedisKeyList() (resp *types.RedisKeyListResp, err error) {
	cachekeys := []types.RedisKey{}
	var cursor uint64 = 0
	for {
		tmpCacheKeys, a, b := l.svcCtx.Redis.ScanCtx(l.ctx, cursor, "cache:verificationSystem:*", 0)
		//fmt.Println("cache a", a)
		//fmt.Println("cache b", b)
		if b != nil {
			break
		}
		cursor = a
		for _, key := range tmpCacheKeys {
			cachekeys = append(cachekeys, types.RedisKey{
				Key: key,
			})
		}
		if a == 0 {
			break
		}
	}

	sort.Slice(cachekeys, func(i, j int) bool {
		return cachekeys[i].Key < cachekeys[j].Key
	})

	/*
		for _, key := range cachekeys {
			fmt.Println("cache key", key)
		}
	*/

	resp = &types.RedisKeyListResp{
		List: cachekeys,
	}

	return
}
