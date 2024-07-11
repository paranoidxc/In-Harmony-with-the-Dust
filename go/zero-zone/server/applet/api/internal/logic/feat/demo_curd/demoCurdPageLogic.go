package demo_curd

import (
	"context"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/utils"
)

type DemoCurdPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoCurdPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoCurdPageLogic {
	return &DemoCurdPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoCurdPageLogic) DemoCurdPage(req *types.DemoCurdPageReq) (resp *types.DemoCurdPageResp, err error) {
	where := " deleted_at IS NULL "
	/*
	   if len(strings.TrimSpace(req.FirmName)) > 0 {
	       where = where + fmt.Sprintf(" AND firm_name LIKE '%s'", "%"+strings.TrimSpace(req.FirmName)+"%")
	   }
	   if len(strings.TrimSpace(req.FirmAlias)) > 0 {
	       where = where + fmt.Sprintf(" AND firm_alias LIKE '%s'", "%"+strings.TrimSpace(req.FirmAlias)+"%")
	   }
	   if len(strings.TrimSpace(req.FirmCode)) > 0 {
	       where = where + fmt.Sprintf(" AND firm_code LIKE '%s'", "%"+strings.TrimSpace(req.FirmCode)+"%")
	   }
	   if len(strings.TrimSpace(req.FirmDesc)) > 0 {
	       where = where + fmt.Sprintf(" AND firm_desc LIKE '%s'", "%"+strings.TrimSpace(req.FirmDesc)+"%")
	   }
	   if len(strings.TrimSpace(req.CreatedAt)) > 0 {
	       where = where + fmt.Sprintf(" AND create_at LIKE '%s'", "%"+strings.TrimSpace(req.CreatedAt)+"%")
	   }
	   if len(strings.TrimSpace(req.UpdatedAt)) > 0 {
	       where = where + fmt.Sprintf(" AND update_at LIKE '%s'", "%"+strings.TrimSpace(req.UpdatedAt)+"%")
	   }
	   if len(strings.TrimSpace(req.DeletedAt)) > 0 {
	       where = where + fmt.Sprintf(" AND delete_at LIKE '%s'", "%"+strings.TrimSpace(req.DeletedAt)+"%")
	   }
	*/

	featDemoCurdPage, err := l.svcCtx.FeatDemoCurdModel.FindPageByWhere(l.ctx, where, req.Page, req.Limit)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var item types.DemoCurd
	DemoCurdPage := make([]types.DemoCurd, 0)
	for _, v := range featDemoCurdPage {
		err := copier.Copy(&item, &v)
		item.CreatedAt = utils.Time2Str(v.CreatedAt)
		item.UpdatedAt = utils.Time2Str(v.UpdatedAt)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		DemoCurdPage = append(DemoCurdPage, item)
	}

	total, err := l.svcCtx.FeatDemoCurdModel.FindPageByWhereCount(l.ctx, where)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	pagination := types.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	return &types.DemoCurdPageResp{
		List:       DemoCurdPage,
		Pagination: pagination,
	}, nil
}
