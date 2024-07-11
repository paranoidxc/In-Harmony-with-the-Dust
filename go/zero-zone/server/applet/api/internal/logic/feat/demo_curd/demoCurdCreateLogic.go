package demo_curd

import (
	"context"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/model"
	errorx2 "zero-zone/pkg/errorx"
)

type DemoCurdCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoCurdCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoCurdCreateLogic {
	return &DemoCurdCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoCurdCreateLogic) DemoCurdCreate(req *types.DemoCurdCreateReq) (err error) {
	var modelParams = new(model.DemoCurd)
	err = copier.Copy(modelParams, req)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}
	_, err = l.svcCtx.FeatDemoCurdModel.Insert(l.ctx, modelParams)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return
}
