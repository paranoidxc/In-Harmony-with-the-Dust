package demo_curd

import (
	"context"
	"strings"

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
	where := []string{"1"}
	if req != nil {
		if req.IncludeDeleted == 0 {
			where = append(where, "t.deleted_at = '0000-00-00 00:00:00'")
		}
	}
	/*
	   if len(strings.TrimSpace(req.FirmName)) > 0 {
	       where = append(where, fmt.Sprintf("firm_name LIKE '%s'", "%"+strings.TrimSpace(req.FirmName)+"%"))
	   }
	   if len(strings.TrimSpace(req.FirmAlias)) > 0 {
	       where = append(where, fmt.Sprintf(" LIKE '%s'", "%"+strings.TrimSpace(req.FirmAlias)+"%"))
	   }
	   if len(strings.TrimSpace(req.FirmCode)) > 0 {
	       where = append(where, fmt.Sprintf("firm_code LIKE '%s'", "%"+strings.TrimSpace(req.FirmCode)+"%"))
	   }
	   if len(strings.TrimSpace(req.FirmDesc)) > 0 {
	       where = append(where, fmt.Sprintf(" LIKE '%s'", "%"+strings.TrimSpace(req.FirmDesc)+"%"))
	   }
	   if len(strings.TrimSpace(req.CreatedAt)) > 0 {
	       where = append(where, fmt.Sprintf(" LIKE '%s'", "%"+strings.TrimSpace(req.CreatedAt)+"%"))
	   }
	   if len(strings.TrimSpace(req.UpdatedAt)) > 0 {
	       where = append(where, fmt.Sprintf(" LIKE '%s'", "%"+strings.TrimSpace(req.UpdatedAt)+"%"))
	   }
	   if len(strings.TrimSpace(req.DeletedAt)) > 0 {
	       where = append(where, fmt.Sprintf(" LIKE '%s'", "%"+strings.TrimSpace(req.DeletedAt)+"%"))
	   }
	*/

	whereStr := strings.Join(where, " AND ")
	featDemoCurdPage, err := l.svcCtx.FeatDemoCurdModel.FindPageByWhere(l.ctx, whereStr, req.Page, req.Limit)
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

	total, err := l.svcCtx.FeatDemoCurdModel.FindPageByWhereCount(l.ctx, whereStr)
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
