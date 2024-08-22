package sys_region

import (
	"context"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	
	"github.com/zeromicro/go-zero/core/logx"
	
	errorx2 "zero-zone/pkg/errorx"
	
)

type SysRegionDeletesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSysRegionDeletesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SysRegionDeletesLogic {
	return &SysRegionDeletesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SysRegionDeletesLogic) SysRegionDeletes(req *types.SysRegionDeletesReq) (err error) {
	if len(req.Id) > 0  {
		err = l.svcCtx.FeatSysRegionModel.Deletes(l.ctx, req.Id)
		if err != nil {
			return  errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
	} else {
		return errorx2.NewSystemError(errorx2.ParamErrorCode, err.Error())
	}

	return
}

