package menu

import (
	"context"
	"encoding/json"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"
	"zero-zone/pkg/utils"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
	"zero-zone/applet/model"
)

type UpdateSysPermMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysPermMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysPermMenuLogic {
	return &UpdateSysPermMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysPermMenuLogic) UpdateSysPermMenu(req *types.UpdateSysPermMenuReq) error {
	if req.ParentId != globalkey.SysTopParentId {
		parentPermMenu, err := l.svcCtx.SysPermMenuModel.FindOne(l.ctx, req.ParentId)
		if err != nil {
			return errorx2.NewDefaultError(errorx2.ParentPermMenuIdErrorCode)
		}

		if parentPermMenu.Type == 2 {
			return errorx2.NewDefaultError(errorx2.SetParentTypeErrorCode)
		}
	}

	/*
		if req.Id <= globalkey.SysProtectPermMenuMaxId {
			return errorx2.NewDefaultError(errorx2.ForbiddenErrorCode)
		}
	*/

	if req.Id == req.ParentId {
		return errorx2.NewDefaultError(errorx2.ParentPermMenuErrorCode)
	}

	permMenuIds := make([]int64, 0)
	permMenuIds = l.getSubPermMenu(permMenuIds, req.Id)
	if utils.ArrayContainValue(permMenuIds, req.ParentId) {
		return errorx2.NewDefaultError(errorx2.SetParentIdErrorCode)
	}

	permMenu, err := l.svcCtx.SysPermMenuModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewDefaultError(errorx2.PermMenuIdErrorCode)
	}

	err = copier.Copy(permMenu, req)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	bytes, err := json.Marshal(req.Perms)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	permMenu.Perms = string(bytes)
	err = l.svcCtx.SysPermMenuModel.Update(l.ctx, permMenu)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return nil
}

func (l *UpdateSysPermMenuLogic) getSubPermMenu(permMenuIds []int64, id int64) []int64 {
	permMenuList, err := l.svcCtx.SysPermMenuModel.FindSubPermMenu(l.ctx, id)
	if err != nil && err != model.ErrNotFound {
		return permMenuIds
	}

	for _, v := range permMenuList {
		permMenuIds = append(permMenuIds, v.Id)
		permMenuIds = l.getSubPermMenu(permMenuIds, v.Id)
	}

	return permMenuIds
}
