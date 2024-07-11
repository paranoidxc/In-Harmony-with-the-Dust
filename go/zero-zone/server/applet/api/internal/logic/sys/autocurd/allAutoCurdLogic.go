package autocurd

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
	"zero-zone/applet/model"
)

type AllAutoCurdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增
func NewAllAutoCurdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllAutoCurdLogic {
	return &AllAutoCurdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllAutoCurdLogic) AllAutoCurd() (resp *types.AutoCurdListResp, err error) {
	// 获取结构体

	List := make([]types.AutoCurd, 0)
	for k, _ := range model.AutoCrudModelList {
		item := types.AutoCurd{
			Name: k,
		}
		List = append(List, item)
	}

	return &types.AutoCurdListResp{
		List: List,
	}, nil
}
