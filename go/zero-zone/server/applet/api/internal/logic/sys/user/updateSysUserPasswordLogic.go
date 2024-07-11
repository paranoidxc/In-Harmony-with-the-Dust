package user

import (
	"context"
	"time"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/utils"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type UpdateSysUserPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysUserPasswordLogic {
	return &UpdateSysUserPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysUserPasswordLogic) UpdateSysUserPassword(req *types.UpdateSysUserPasswordReq) error {
	sysUser, err := l.svcCtx.SysUserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewDefaultError(errorx2.UserIdErrorCode)
	}

	err = copier.Copy(sysUser, req)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	sysUser.Password = utils.MD5(utils.MD5(time.Now().Format("20060102")) + l.svcCtx.Config.Salt)
	//sysUser.Password = utils.MD5(req.Password + l.svcCtx.Config.Salt)
	err = l.svcCtx.SysUserModel.Update(l.ctx, sysUser)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return nil
}
