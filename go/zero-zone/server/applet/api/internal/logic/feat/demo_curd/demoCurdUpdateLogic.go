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

type DemoCurdUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoCurdUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoCurdUpdateLogic {
	return &DemoCurdUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoCurdUpdateLogic) DemoCurdUpdate(req *types.DemoCurdUpdateReq) (err error) {
	modelParams := &model.DemoCurd{}
	modelParams, err = l.svcCtx.FeatDemoCurdModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewDefaultError(errorx2.UserIdErrorCode)
	}

	err = copier.Copy(modelParams, req)
	if err != nil {
		logx.Error("复制参数失败", err)
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	err = l.svcCtx.FeatDemoCurdModel.Update(l.ctx, modelParams)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return
}
