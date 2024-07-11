package dict

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type UpdateConfigDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateConfigDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateConfigDictLogic {
	return &UpdateConfigDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateConfigDictLogic) UpdateConfigDict(req *types.UpdateConfigDictReq) error {
	if req.ParentId != globalkey.SysTopParentId {
		_, err := l.svcCtx.SysDictionaryModel.FindOne(l.ctx, req.ParentId)
		if err != nil {
			return errorx2.NewDefaultError(errorx2.ParentDictionaryIdErrorCode)
		}
	}

	configDictionary, err := l.svcCtx.SysDictionaryModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewDefaultError(errorx2.DictionaryIdErrorCode)
	}

	err = copier.Copy(configDictionary, req)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	err = l.svcCtx.SysDictionaryModel.Update(l.ctx, configDictionary)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return nil
}
