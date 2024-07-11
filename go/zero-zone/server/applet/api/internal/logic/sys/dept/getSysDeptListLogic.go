package dept

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type GetSysDeptListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSysDeptListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSysDeptListLogic {
	return &GetSysDeptListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSysDeptListLogic) GetSysDeptList() (resp *types.SysDeptListResp, err error) {
	sysDeptList, err := l.svcCtx.SysDeptModel.FindAll(l.ctx)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var dept types.Dept
	deptList := make([]types.Dept, 0)
	for _, v := range sysDeptList {
		err := copier.Copy(&dept, &v)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		deptList = append(deptList, dept)
	}

	return &types.SysDeptListResp{
		List: deptList,
	}, nil
}

func (l *GetSysDeptListLogic) GetSysDeptPage(req *types.SysDeptPageReq) (resp *types.SysDeptPageResp, err error) {
	sysDeptList, err := l.svcCtx.SysDeptModel.FindPage(l.ctx, req.Page, req.Limit)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var dept types.Dept
	deptList := make([]types.Dept, 0)
	for _, sysDept := range sysDeptList {
		err := copier.Copy(&dept, &sysDept)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		deptList = append(deptList, dept)
	}

	total, err := l.svcCtx.SysDeptModel.FindCount(l.ctx)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	pagination := types.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	return &types.SysDeptPageResp{
		List:       deptList,
		Pagination: pagination,
	}, nil
}
