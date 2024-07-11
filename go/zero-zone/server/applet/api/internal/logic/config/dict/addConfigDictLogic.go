package dict

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

type AddConfigDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddConfigDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddConfigDictLogic {
	return &AddConfigDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddConfigDictLogic) AddConfigDict(req *types.AddConfigDictReq) error {
	if req.ParentId != globalkey.SysTopParentId {
		_, err := l.svcCtx.SysDictionaryModel.FindOne(l.ctx, req.ParentId)
		if err != nil {
			return errorx2.NewDefaultError(errorx2.ParentDictionaryIdErrorCode)
		}
	}
	_, err := l.svcCtx.SysDictionaryModel.FindOneByUniqueKey(l.ctx, req.UniqueKey)
	if err == model.ErrNotFound {
		var dictionary = new(model.SysDictionary)
		err = copier.Copy(dictionary, req)
		if err != nil {
			return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		_, err = l.svcCtx.SysDictionaryModel.Insert(l.ctx, dictionary)
		if err != nil {
			return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}

		return nil
	} else {

		return errorx2.NewDefaultError(errorx2.AddDictionaryErrorCode)
	}
}
