package sys_region

import (
	"context"
	"zero-zone/pkg/utils"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/model"
)

type SysRegionDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSysRegionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SysRegionDetailLogic {
	return &SysRegionDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SysRegionDetailLogic) SysRegionDetail(req *types.SysRegionDetailReq) (resp *types.SysRegionDetailResp, err error) {
	resp = &types.SysRegionDetailResp{}
	item := &model.SysRegion{}
	item, err = l.svcCtx.FeatSysRegionModel.FindOne(l.ctx, req.Id)
	err = copier.Copy(resp, item)
	resp.CreatedAt = utils.Time2Str(item.CreatedAt)
	resp.UpdatedAt = utils.Time2Str(item.UpdatedAt)
	if err != nil {
		logx.Error("复制结果失败", err)
		return nil, err
	}
	return
}
