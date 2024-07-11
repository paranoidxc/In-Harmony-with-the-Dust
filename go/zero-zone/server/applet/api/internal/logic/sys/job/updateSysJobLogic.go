package job

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type UpdateSysJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysJobLogic {
	return &UpdateSysJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysJobLogic) UpdateSysJob(req *types.UpdateSysJobReq) error {
	sysJob, err := l.svcCtx.SysJobModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewDefaultError(errorx2.JobIdErrorCode)
	}

	if req.Status == globalkey.SysDisable {
		count, _ := l.svcCtx.SysUserModel.FindCountByJobId(l.ctx, req.Id)
		if count > 0 {
			return errorx2.NewDefaultError(errorx2.JobIsUsingErrorCode)
		}
	}

	err = copier.Copy(sysJob, req)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	err = l.svcCtx.SysJobModel.Update(l.ctx, sysJob)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return nil
}
