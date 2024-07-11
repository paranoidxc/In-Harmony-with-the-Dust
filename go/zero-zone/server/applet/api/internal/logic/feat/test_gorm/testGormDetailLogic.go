package test_gorm

import (
	"context"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/model"
)

type TestGormDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTestGormDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestGormDetailLogic {
	return &TestGormDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TestGormDetailLogic) TestGormDetail(req *types.TestGormDetailReq) (resp *types.TestGormDetailResp, err error) {
	resp = &types.TestGormDetailResp{}
	item := &model.TestGorm{}
	item, err = l.svcCtx.FeatTestGormModel.FindOne(req.ID)
	err = copier.Copy(resp, item)
	if err != nil {
		logx.Error("复制结果失败", err)
		return nil, err
	}
	return
}
