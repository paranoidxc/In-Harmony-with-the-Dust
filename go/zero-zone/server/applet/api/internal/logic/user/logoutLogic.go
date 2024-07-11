package user

import (
	"context"
	"strconv"
	"zero-zone/pkg/globalkey"
	"zero-zone/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout() error {
	userId := strconv.FormatInt(utils.GetUserId(l.ctx), 10)
	_, _ = l.svcCtx.Redis.Del(globalkey.SysPermMenuCachePrefix + userId)
	_, _ = l.svcCtx.Redis.Del(globalkey.SysOnlineUserCachePrefix + userId)
	_, _ = l.svcCtx.Redis.Del(globalkey.SysUserIdCachePrefix + userId)

	return nil
}
