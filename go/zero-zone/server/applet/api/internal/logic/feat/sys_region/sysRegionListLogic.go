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

type SysRegionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSysRegionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SysRegionListLogic {
	return &SysRegionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SysRegionListLogic) SysRegionList(req *types.SysRegionListReq) (resp *types.SysRegionListResp, err error) {
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
	featSysRegionList, err := l.svcCtx.FeatSysRegionModel.FindAllByWhere(l.ctx, whereStr)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var item types.SysRegion
	SysRegionList := make([]types.SysRegion, 0)
	for _, v := range featSysRegionList {
		err := copier.Copy(&item, &v)
		item.CreatedAt = utils.Time2Str(v.CreatedAt)
		item.UpdatedAt = utils.Time2Str(v.UpdatedAt)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		SysRegionList = append(SysRegionList, item)
	}

	return &types.SysRegionListResp{
		List: SysRegionList,
	}, nil
}
