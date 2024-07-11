package dict

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type GetConfigDictListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConfigDictListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConfigDictListLogic {
	return &GetConfigDictListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConfigDictListLogic) GetConfigDictList() (resp *types.ConfigDictListResp, err error) {
	configDictionaryList, err := l.svcCtx.SysDictionaryModel.FindDictionaryList(l.ctx)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var dictionary types.ConfigDict
	dictionaryList := make([]types.ConfigDict, 0)
	for _, v := range configDictionaryList {
		err := copier.Copy(&dictionary, &v)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		dictionaryList = append(dictionaryList, dictionary)
	}

	return &types.ConfigDictListResp{
		List: dictionaryList,
	}, nil
}
