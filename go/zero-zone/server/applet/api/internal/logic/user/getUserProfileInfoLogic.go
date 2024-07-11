package user

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type GetUserProfileInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserProfileInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserProfileInfoLogic {
	return &GetUserProfileInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserProfileInfoLogic) GetUserProfileInfo() (resp *types.UserProfileInfoResp, err error) {
	userId := utils.GetUserId(l.ctx)
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return &types.UserProfileInfoResp{
		Username: user.Username,
		Nickname: user.Nickname,
		Gender:   user.Gender,
		Email:    user.Email,
		Mobile:   user.Mobile,
		Remark:   user.Remark,
		Avatar:   user.Avatar,
	}, nil
}
