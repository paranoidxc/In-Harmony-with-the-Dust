package test_gorm

import (
	"context"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"

	errorx2 "zero-zone/pkg/errorx"
)

type TestGormDeletesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTestGormDeletesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestGormDeletesLogic {
	return &TestGormDeletesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TestGormDeletesLogic) TestGormDeletes(req *types.TestGormDeletesReq) (err error) {
	if len(req.ID) > 0 {
		err = l.svcCtx.FeatTestGormModel.Deletes(req.ID)
		if err != nil {
			return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
	} else {
		return errorx2.NewSystemError(errorx2.ParamErrorCode, err.Error())
	}

	return
}
