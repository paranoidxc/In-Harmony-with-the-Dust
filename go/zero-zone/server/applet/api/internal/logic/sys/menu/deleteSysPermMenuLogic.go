package menu

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"

	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type DeleteSysPermMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSysPermMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSysPermMenuLogic {
	return &DeleteSysPermMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSysPermMenuLogic) DeleteSysPermMenu(req *types.DeleteSysPermMenuReq) error {
	if req.Id <= globalkey.SysProtectPermMenuMaxId {
		return errorx2.NewDefaultError(errorx2.ForbiddenErrorCode)
	}

	count, _ := l.svcCtx.SysPermMenuModel.FindCountByParentId(l.ctx, req.Id)
	if count != 0 {
		return errorx2.NewDefaultError(errorx2.DeletePermMenuErrorCode)
	}

	err := l.svcCtx.SysPermMenuModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return nil
}
