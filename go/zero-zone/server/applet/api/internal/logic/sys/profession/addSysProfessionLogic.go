package profession

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
	"zero-zone/applet/model"
)

type AddSysProfessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSysProfessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSysProfessionLogic {
	return &AddSysProfessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSysProfessionLogic) AddSysProfession(req *types.AddSysProfessionReq) error {
	_, err := l.svcCtx.SysProfessionModel.FindOneByName(l.ctx, req.Name)
	if err == model.ErrNotFound {
		var sysProfession = new(model.SysProfession)
		err = copier.Copy(sysProfession, req)
		if err != nil {
			return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		_, err = l.svcCtx.SysProfessionModel.Insert(l.ctx, sysProfession)
		if err != nil {
			return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}

		return nil
	} else {

		return errorx2.NewDefaultError(errorx2.AddProfessionErrorCode)
	}
}
