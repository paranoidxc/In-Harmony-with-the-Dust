package role

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"

	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type DeleteSysRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSysRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSysRoleLogic {
	return &DeleteSysRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSysRoleLogic) DeleteSysRole(req *types.DeleteSysRoleReq) error {
	if req.Id == globalkey.SysSuperRoleId {
		return errorx2.NewDefaultError(errorx2.ForbiddenErrorCode)
	}

	roleList, _ := l.svcCtx.SysRoleModel.FindSubRole(l.ctx, req.Id)
	if len(roleList) != 0 {
		return errorx2.NewDefaultError(errorx2.DeleteRoleErrorCode)
	}

	count, _ := l.svcCtx.SysUserModel.FindCountByRoleId(l.ctx, req.Id)
	if count != 0 {
		return errorx2.NewDefaultError(errorx2.RoleIsUsingErrorCode)
	}

	err := l.svcCtx.SysRoleModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return nil
}
