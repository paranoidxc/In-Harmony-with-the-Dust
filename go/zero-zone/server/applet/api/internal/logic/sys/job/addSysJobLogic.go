package job

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
	"zero-zone/applet/model"
)

type AddSysJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSysJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSysJobLogic {
	return &AddSysJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSysJobLogic) AddSysJob(req *types.AddSysJobReq) error {
	_, err := l.svcCtx.SysJobModel.FindOneByName(l.ctx, req.Name)
	if err == model.ErrNotFound {
		var sysJob = new(model.SysJob)
		err = copier.Copy(sysJob, req)
		if err != nil {
			return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		_, err = l.svcCtx.SysJobModel.Insert(l.ctx, sysJob)
		if err != nil {
			return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}

		return nil
	} else {

		return errorx2.NewDefaultError(errorx2.AddJobErrorCode)
	}
}
