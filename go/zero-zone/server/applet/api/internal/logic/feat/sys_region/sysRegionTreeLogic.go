package sys_region

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
	"zero-zone/applet/model"
	errorx2 "zero-zone/pkg/errorx"
)

type SysRegionTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSysRegionTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SysRegionTreeLogic {
	return &SysRegionTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SysRegionTreeLogic) SysRegionTree() (resp *types.SysRegionTreeResp, err error) {
	where := []string{"1"}
	whereStr := strings.Join(where, " AND ")
	featSysRegionList, err := l.svcCtx.FeatSysRegionModel.FindAllByWhere(l.ctx, whereStr)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	treeData := buildNestedJSON(featSysRegionList, 0)
	//logc.Infow(l.ctx, "buildNestedJSON", logc.Field("rer", treeData))

	resp = &types.SysRegionTreeResp{
		TreeData: treeData,
	}

	return
}

func buildNestedJSON(regions []*model.SysRegion, parentNo int64) []types.SysRegionTree {
	ItemData := []types.SysRegionTree{}
	for _, region := range regions {
		if region.ParentNo == parentNo {
			childNode := buildNestedJSON(regions, region.No)
			node := types.SysRegionTree{
				Name:     region.Name,
				No:       region.No,
				PySzm:    region.PySzm,
				Children: childNode,
			}
			ItemData = append(ItemData, node)
		}
	}
	return ItemData
}
