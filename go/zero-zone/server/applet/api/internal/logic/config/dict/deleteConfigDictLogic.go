package dict

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"

	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type DeleteConfigDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteConfigDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteConfigDictLogic {
	return &DeleteConfigDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteConfigDictLogic) DeleteConfigDict(req *types.DeleteConfigDictReq) error {
	if req.Id <= globalkey.SysProtectDictionaryMaxId {
		return errorx2.NewDefaultError(errorx2.ForbiddenErrorCode)
	}

	total, err := l.svcCtx.SysDictionaryModel.FindCountByParentId(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	if total > 0 {
		return errorx2.NewDefaultError(errorx2.DeleteDictionaryErrorCode)
	}

	err = l.svcCtx.SysDictionaryModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return nil
}
