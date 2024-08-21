package wechat

import (
	"context"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/app/model"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/utils"
)

type WechatLoginSucLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWechatLoginSucLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WechatLoginSucLogic {
	return &WechatLoginSucLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WechatLoginSucLogic) WechatLoginSuc(req *types.WechatLoginReq) (resp *types.WechatLoginResq, err error) {
	// todo: add your logic here and delete this line

	return
}
