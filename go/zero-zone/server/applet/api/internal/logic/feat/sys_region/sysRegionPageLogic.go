package sys_region

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

type SysRegionPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSysRegionPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SysRegionPageLogic {
	return &SysRegionPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SysRegionPageLogic) SysRegionPage(req *types.SysRegionPageReq) (resp *types.SysRegionPageResp, err error) {
	where := []string{"1"}
	if req != nil {
		if req.IncludeDeleted == 0 {
			where = append(where, "t.deleted_at = '0000-00-00 00:00:00'")
		}
	}
	/*
	   if len(strings.TrimSpace(req.No)) > 0 {
	       where = append(where, fmt.Sprintf("no LIKE '%s'", "%"+strings.TrimSpace(req.No)+"%"))
	   }
	   if len(strings.TrimSpace(req.Name)) > 0 {
	       where = append(where, fmt.Sprintf("no LIKE '%s'", "%"+strings.TrimSpace(req.Name)+"%"))
	   }
	   if len(strings.TrimSpace(req.ParentNo)) > 0 {
	       where = append(where, fmt.Sprintf("no LIKE '%s'", "%"+strings.TrimSpace(req.ParentNo)+"%"))
	   }
	   if len(strings.TrimSpace(req.Code)) > 0 {
	       where = append(where, fmt.Sprintf("no LIKE '%s'", "%"+strings.TrimSpace(req.Code)+"%"))
	   }
	   if len(strings.TrimSpace(req.TypeName)) > 0 {
	       where = append(where, fmt.Sprintf("no LIKE '%s'", "%"+strings.TrimSpace(req.TypeName)+"%"))
	   }
	   if len(strings.TrimSpace(req.PySzm)) > 0 {
	       where = append(where, fmt.Sprintf("no LIKE '%s'", "%"+strings.TrimSpace(req.PySzm)+"%"))
	   }
	   if len(strings.TrimSpace(req.CreatedAt)) > 0 {
	       where = append(where, fmt.Sprintf("created_at LIKE '%s'", "%"+strings.TrimSpace(req.CreatedAt)+"%"))
	   }
	   if len(strings.TrimSpace(req.UpdatedAt)) > 0 {
	       where = append(where, fmt.Sprintf("updated_at LIKE '%s'", "%"+strings.TrimSpace(req.UpdatedAt)+"%"))
	   }
	   if len(strings.TrimSpace(req.DeletedAt)) > 0 {
	       where = append(where, fmt.Sprintf("deleted_at LIKE '%s'", "%"+strings.TrimSpace(req.DeletedAt)+"%"))
	   }
	   if len(strings.TrimSpace(req.IsDel)) > 0 {
	       where = append(where, fmt.Sprintf("deleted_at LIKE '%s'", "%"+strings.TrimSpace(req.IsDel)+"%"))
	   }
	*/

	whereStr := strings.Join(where, " AND ")
	featSysRegionPage, err := l.svcCtx.FeatSysRegionModel.FindPageByWhere(l.ctx, whereStr, req.Page, req.Limit)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var item types.SysRegion
	SysRegionPage := make([]types.SysRegion, 0)
	for _, v := range featSysRegionPage {
		err := copier.Copy(&item, &v)
		item.No = utils.Int642Str(v.No)
		item.ParentNo = utils.Int642Str(v.ParentNo)
		item.TypeName = v.Typename
		item.CreatedAt = utils.Time2Str(v.CreatedAt)
		item.UpdatedAt = utils.Time2Str(v.UpdatedAt)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		SysRegionPage = append(SysRegionPage, item)
	}

	total, err := l.svcCtx.FeatSysRegionModel.FindPageByWhereCount(l.ctx, whereStr)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	pagination := types.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	return &types.SysRegionPageResp{
		List:       SysRegionPage,
		Pagination: pagination,
	}, nil
}
