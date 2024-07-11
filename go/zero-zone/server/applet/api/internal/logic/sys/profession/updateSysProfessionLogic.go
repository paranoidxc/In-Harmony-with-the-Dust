package profession

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type UpdateSysProfessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysProfessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysProfessionLogic {
	return &UpdateSysProfessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysProfessionLogic) UpdateSysProfession(req *types.UpdateSysProfessionReq) error {
	sysProfession, err := l.svcCtx.SysProfessionModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewDefaultError(errorx2.ProfessionIdErrorCode)
	}

	if req.Status == globalkey.SysDisable {
		count, _ := l.svcCtx.SysUserModel.FindCountByProfessionId(l.ctx, req.Id)
		if count > 0 {
			return errorx2.NewDefaultError(errorx2.JobIsUsingErrorCode)
		}
	}

	err = copier.Copy(sysProfession, req)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	err = l.svcCtx.SysProfessionModel.Update(l.ctx, sysProfession)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return nil
}
