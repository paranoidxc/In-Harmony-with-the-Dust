package demo_curd

import (
	"context"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	
	"github.com/zeromicro/go-zero/core/logx"
	
	errorx2 "zero-zone/pkg/errorx"
	
)

type DemoCurdDeletesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoCurdDeletesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoCurdDeletesLogic {
	return &DemoCurdDeletesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoCurdDeletesLogic) DemoCurdDeletes(req *types.DemoCurdDeletesReq) (err error) {
	if len(req.Id) > 0  {
		err = l.svcCtx.FeatDemoCurdModel.Deletes(l.ctx, req.Id)
		if err != nil {
			return  errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
	} else {
		return errorx2.NewSystemError(errorx2.ParamErrorCode, err.Error())
	}

	return
}

