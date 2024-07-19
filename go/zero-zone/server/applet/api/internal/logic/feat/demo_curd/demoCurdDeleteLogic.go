package demo_curd

import (
	"context"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	
	"github.com/zeromicro/go-zero/core/logx"
	
	errorx2 "zero-zone/pkg/errorx"
	
)

type DemoCurdDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoCurdDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoCurdDeleteLogic {
	return &DemoCurdDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoCurdDeleteLogic) DemoCurdDelete(req *types.DemoCurdDeleteReq) (err error) {
	err = l.svcCtx.FeatDemoCurdModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return
}

