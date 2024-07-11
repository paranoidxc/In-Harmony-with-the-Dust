package profession

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type DeleteSysProfessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSysProfessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSysProfessionLogic {
	return &DeleteSysProfessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSysProfessionLogic) DeleteSysProfession(req *types.DeleteSysProfessionReq) error {
	count, _ := l.svcCtx.SysUserModel.FindCountByCondition(l.ctx, "profession_id", req.Id)
	if count != 0 {
		return errorx2.NewDefaultError(errorx2.DeleteProfessionErrorCode)
	}

	err := l.svcCtx.SysProfessionModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return nil
}
