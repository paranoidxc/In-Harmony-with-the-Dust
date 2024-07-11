package demo_curd

import (
	"context"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/model"

	"zero-zone/pkg/utils"
)

type DemoCurdDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoCurdDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoCurdDetailLogic {
	return &DemoCurdDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoCurdDetailLogic) DemoCurdDetail(req *types.DemoCurdDetailReq) (resp *types.DemoCurdDetailResp, err error) {
	resp = &types.DemoCurdDetailResp{}
	item := &model.DemoCurd{}
	item, err = l.svcCtx.FeatDemoCurdModel.FindOne(l.ctx, req.Id)
	err = copier.Copy(resp, item)
	resp.CreatedAt = utils.Time2Str(item.CreatedAt)
	resp.UpdatedAt = utils.Time2Str(item.UpdatedAt)
	if err != nil {
		logx.Error("复制结果失败", err)
		return nil, err
	}
	return
}
