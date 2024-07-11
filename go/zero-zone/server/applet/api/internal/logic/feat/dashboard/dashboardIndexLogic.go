package dashboard

import (
	"context"
	"fmt"
	"strings"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	//"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	//"zero-zone/app/model"
	//errorx2 "zero-zone/pkg/errorx"
	//"zero-zone/pkg/utils"
)

type DashboardIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDashboardIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DashboardIndexLogic {
	return &DashboardIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DashboardIndexLogic) DashboardIndex(req *types.DashboardIndexReq) (resp *types.DashboardIndexResp, err error) {
	where := " t.deleted_at IS NULL AND t.status = 1"
	if len(strings.TrimSpace(req.DateRangeStart)) > 0 {
		where = where + fmt.Sprintf(" AND t.updated_at >= '%s 00:00:00'", strings.TrimSpace(req.DateRangeStart))
	}
	if len(strings.TrimSpace(req.DateRangeEnd)) > 0 {
		where = where + fmt.Sprintf(" AND t.updated_at <= '%s 23:59:59'", strings.TrimSpace(req.DateRangeEnd))
	}

	resp = &types.DashboardIndexResp{
		HxOrderCnt:  0,
		UhxOrderCnt: 0,
	}

	return
}
