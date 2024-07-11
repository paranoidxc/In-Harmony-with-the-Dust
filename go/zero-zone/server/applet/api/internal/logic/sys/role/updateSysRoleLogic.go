package role

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

type UpdateSysRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysRoleLogic {
	return &UpdateSysRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysRoleLogic) UpdateSysRole(req *types.UpdateSysRoleReq) error {
	if req.ParentId != globalkey.SysTopParentId {
		_, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, req.ParentId)
		if err != nil {
			return errorx2.NewDefaultError(errorx2.ParentRoleIdErrorCode)
		}
	}

	if req.Id == globalkey.SysSuperRoleId {
		return errorx2.NewDefaultError(errorx2.NotPermMenuErrorCode)
	}

	if req.Id == req.ParentId {
		return errorx2.NewDefaultError(errorx2.ParentRoleErrorCode)
	}

	role, err := l.svcCtx.SysRoleModel.FindOneByUniqueKey(l.ctx, req.UniqueKey)
	if err != model.ErrNotFound && role.Id != req.Id {
		return errorx2.NewDefaultError(errorx2.UpdateRoleUniqueKeyErrorCode)
	}

	roleIds := make([]int64, 0)
	roleIds = l.getSubRole(roleIds, req.Id)
	if utils.ArrayContainValue(roleIds, req.ParentId) {
		return errorx2.NewDefaultError(errorx2.SetParentIdErrorCode)
	}

	sysRole, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewDefaultError(errorx2.RoleIdErrorCode)
	}

	err = copier.Copy(sysRole, req)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	bytes, err := json.Marshal(req.PermMenuIds)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	bytes2, err := json.Marshal(req.PermMenuIdsAll)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	sysRole.PermMenuIds = string(bytes)
	sysRole.PermMenuIdsAll = string(bytes2)
	err = l.svcCtx.SysRoleModel.Update(l.ctx, sysRole)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return nil
}

func (l *UpdateSysRoleLogic) getSubRole(roleIds []int64, id int64) []int64 {
	roleList, err := l.svcCtx.SysRoleModel.FindSubRole(l.ctx, id)
	if err != nil && err != model.ErrNotFound {
		return roleIds
	}

	for _, v := range roleList {
		roleIds = append(roleIds, v.Id)
		roleIds = l.getSubRole(roleIds, v.Id)
	}

	return roleIds
}
