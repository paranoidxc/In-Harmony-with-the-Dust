package menu

import (
	"context"
	"encoding/json"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
	"zero-zone/applet/model"
)

type AddSysPermMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSysPermMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSysPermMenuLogic {
	return &AddSysPermMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSysPermMenuLogic) AddSysPermMenu(req *types.AddSysPermMenuReq) error {
	if req.ParentId != globalkey.SysTopParentId {
		parentPermMenu, err := l.svcCtx.SysPermMenuModel.FindOne(l.ctx, req.ParentId)
		if err != nil {
			return errorx2.NewDefaultError(errorx2.ParentPermMenuIdErrorCode)
		}

		if parentPermMenu.Type == 2 {
			return errorx2.NewDefaultError(errorx2.SetParentTypeErrorCode)
		}
	}

	var permMenu = new(model.SysPermMenu)
	err := copier.Copy(permMenu, req)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	bytes, err := json.Marshal(req.Perms)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	permMenu.Perms = string(bytes)
	_, err = l.svcCtx.SysPermMenuModel.Insert(l.ctx, permMenu)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return nil
}
