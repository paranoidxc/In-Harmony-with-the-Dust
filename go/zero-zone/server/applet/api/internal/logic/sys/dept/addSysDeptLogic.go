package dept

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
	"zero-zone/applet/model"
)

type AddSysDeptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSysDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSysDeptLogic {
	return &AddSysDeptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSysDeptLogic) AddSysDept(req *types.AddSysDeptReq) error {
	_, err := l.svcCtx.SysDeptModel.FindOneByUniqueKey(l.ctx, req.UniqueKey)
	if err == model.ErrNotFound {
		if req.ParentId != globalkey.SysTopParentId {
			_, err := l.svcCtx.SysDeptModel.FindOne(l.ctx, req.ParentId)
			if err != nil {
				return errorx2.NewDefaultError(errorx2.ParentDeptIdErrorCode)
			}
		}

		var sysDept = new(model.SysDept)
		err = copier.Copy(sysDept, req)
		if err != nil {
			return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		_, err = l.svcCtx.SysDeptModel.Insert(l.ctx, sysDept)
		if err != nil {
			return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		return nil
	} else {
		return errorx2.NewDefaultError(errorx2.AddDeptErrorCode)
	}
}
