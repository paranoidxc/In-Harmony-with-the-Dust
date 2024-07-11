package test_gorm

import (
	"context"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"

	errorx2 "zero-zone/pkg/errorx"
)

type TestGormDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTestGormDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestGormDeleteLogic {
	return &TestGormDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TestGormDeleteLogic) TestGormDelete(req *types.TestGormDeleteReq) (err error) {
	err = l.svcCtx.FeatTestGormModel.Delete(req.ID)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return
}
