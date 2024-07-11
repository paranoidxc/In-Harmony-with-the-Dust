package dept

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type DeleteSysDeptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSysDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSysDeptLogic {
	return &DeleteSysDeptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSysDeptLogic) DeleteSysDept(req *types.DeleteSysDeptReq) error {
	deptList, _ := l.svcCtx.SysDeptModel.FindSubDept(l.ctx, req.Id)
	if len(deptList) != 0 {
		return errorx2.NewDefaultError(errorx2.DeleteDeptErrorCode)
	}

	count, _ := l.svcCtx.SysUserModel.FindCountByCondition(l.ctx, "dept_id", req.Id)
	if count != 0 {
		return errorx2.NewDefaultError(errorx2.DeptHasUserErrorCode)
	}

	err := l.svcCtx.SysDeptModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return nil
}
