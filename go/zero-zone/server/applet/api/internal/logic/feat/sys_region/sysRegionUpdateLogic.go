package sys_region

import (
	"context"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/model"
	errorx2 "zero-zone/pkg/errorx"
	
)

type SysRegionUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSysRegionUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SysRegionUpdateLogic {
	return &SysRegionUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SysRegionUpdateLogic) SysRegionUpdate(req *types.SysRegionUpdateReq) (err error) {
	modelParams := &model.SysRegion{}
	modelParams, err = l.svcCtx.FeatSysRegionModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewDefaultError(errorx2.UserIdErrorCode)
	}

	err = copier.Copy(modelParams, req)
	if err != nil {
		logx.Error("复制参数失败", err)
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	err = l.svcCtx.FeatSysRegionModel.Update(l.ctx, modelParams)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return
}

