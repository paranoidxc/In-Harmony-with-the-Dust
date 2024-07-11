package job

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type DeleteSysJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSysJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSysJobLogic {
	return &DeleteSysJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSysJobLogic) DeleteSysJob(req *types.DeleteSysJobReq) error {
	count, _ := l.svcCtx.SysUserModel.FindCountByCondition(l.ctx, "job_id", req.Id)
	if count != 0 {
		return errorx2.NewDefaultError(errorx2.DeleteJobErrorCode)
	}

	err := l.svcCtx.SysJobModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return nil
}
