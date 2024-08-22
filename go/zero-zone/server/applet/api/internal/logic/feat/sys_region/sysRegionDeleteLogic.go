package sys_region

import (
	"context"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	
	"github.com/zeromicro/go-zero/core/logx"
	
	errorx2 "zero-zone/pkg/errorx"
	
)

type SysRegionDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSysRegionDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SysRegionDeleteLogic {
	return &SysRegionDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SysRegionDeleteLogic) SysRegionDelete(req *types.SysRegionDeleteReq) (err error) {
	err = l.svcCtx.FeatSysRegionModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return
}

