package test_gorm

import (
	"context"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/model"
	errorx2 "zero-zone/pkg/errorx"
)

type TestGormUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTestGormUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestGormUpdateLogic {
	return &TestGormUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TestGormUpdateLogic) TestGormUpdate(req *types.TestGormUpdateReq) (err error) {
	modelParams := &model.TestGorm{}
	modelParams, err = l.svcCtx.FeatTestGormModel.FindOne(req.ID)
	if err != nil {
		return errorx2.NewDefaultError(errorx2.UserIdErrorCode)
	}

	err = copier.Copy(modelParams, req)
	if err != nil {
		logx.Error("复制参数失败", err)
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	err = l.svcCtx.FeatTestGormModel.Update(modelParams)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return
}
