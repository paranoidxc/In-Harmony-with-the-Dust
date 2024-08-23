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

type DemoCurdListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoCurdListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoCurdListLogic {
	return &DemoCurdListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoCurdListLogic) DemoCurdList(req *types.DemoCurdListReq) (resp *types.DemoCurdListResp, err error) {
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
	featDemoCurdList, err := l.svcCtx.FeatDemoCurdModel.FindAllByWhere(l.ctx, whereStr)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var item types.DemoCurd
	DemoCurdList := make([]types.DemoCurd, 0)
	for _, v := range featDemoCurdList {
		err := copier.Copy(&item, &v)
		item.CreatedAt = utils.Time2Str(v.CreatedAt)
		item.UpdatedAt = utils.Time2Str(v.UpdatedAt)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		DemoCurdList = append(DemoCurdList, item)
	}

	return &types.DemoCurdListResp{
		List: DemoCurdList,
	}, nil
}
