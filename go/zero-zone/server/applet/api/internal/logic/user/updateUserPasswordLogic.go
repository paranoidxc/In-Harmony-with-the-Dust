package user

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"
	"zero-zone/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type UpdateUserPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserPasswordLogic {
	return &UpdateUserPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserPasswordLogic) UpdateUserPassword(req *types.UpdatePasswordReq) error {
	dictionary, err := l.svcCtx.SysDictionaryModel.FindOneByUniqueKey(l.ctx, "sys_ch_pwd")
	if dictionary.Status == globalkey.SysDisable {
		return errorx2.NewDefaultError(errorx2.ForbiddenErrorCode)
	}

	userId := utils.GetUserId(l.ctx)
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	if user.Password != utils.MD5(req.OldPassword+l.svcCtx.Config.Salt) {
		return errorx2.NewDefaultError(errorx2.PasswordErrorCode)
	}

	user.Password = utils.MD5(req.NewPassword + l.svcCtx.Config.Salt)
	err = l.svcCtx.SysUserModel.Update(l.ctx, user)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return nil
}
