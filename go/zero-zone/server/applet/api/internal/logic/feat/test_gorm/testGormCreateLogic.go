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

type TestGormCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTestGormCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestGormCreateLogic {
	return &TestGormCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TestGormCreateLogic) TestGormCreate(req *types.TestGormCreateReq) (err error) {
	var modelParams = new(model.TestGorm)
	err = copier.Copy(modelParams, req)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}
	err = l.svcCtx.FeatTestGormModel.Insert(modelParams)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return
}
