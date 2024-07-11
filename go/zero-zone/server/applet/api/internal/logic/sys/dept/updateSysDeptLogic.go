package dept

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"
	"zero-zone/pkg/utils"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
	"zero-zone/applet/model"
)

type UpdateSysDeptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysDeptLogic {
	return &UpdateSysDeptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysDeptLogic) UpdateSysDept(req *types.UpdateSysDeptReq) error {
	if req.ParentId != globalkey.SysTopParentId {
		_, err := l.svcCtx.SysDeptModel.FindOne(l.ctx, req.ParentId)
		if err != nil {
			return errorx2.NewDefaultError(errorx2.ParentDeptIdErrorCode)
		}
	}

	if req.Id == req.ParentId {
		return errorx2.NewDefaultError(errorx2.ParentDeptErrorCode)
	}

	dept, err := l.svcCtx.SysDeptModel.FindOneByUniqueKey(l.ctx, req.UniqueKey)
	if err != model.ErrNotFound && dept.Id != req.Id {
		return errorx2.NewDefaultError(errorx2.UpdateDeptUniqueKeyErrorCode)
	}

	deptIds := make([]int64, 0)
	deptIds = l.getSubDept(deptIds, req.Id)
	if utils.ArrayContainValue(deptIds, req.ParentId) {
		return errorx2.NewDefaultError(errorx2.SetParentIdErrorCode)
	}

	sysDept, err := l.svcCtx.SysDeptModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewDefaultError(errorx2.DeptIdErrorCode)
	}

	err = copier.Copy(sysDept, req)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	err = l.svcCtx.SysDeptModel.Update(l.ctx, sysDept)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return nil
}

func (l *UpdateSysDeptLogic) getSubDept(deptIds []int64, id int64) []int64 {
	deptList, err := l.svcCtx.SysDeptModel.FindSubDept(l.ctx, id)
	if err != nil && err != model.ErrNotFound {
		return deptIds
	}

	for _, v := range deptList {
		deptIds = append(deptIds, v.Id)
		deptIds = l.getSubDept(deptIds, v.Id)
	}

	return deptIds
}
