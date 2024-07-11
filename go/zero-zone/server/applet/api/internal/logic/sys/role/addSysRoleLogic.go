package role

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

type AddSysRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSysRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSysRoleLogic {
	return &AddSysRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSysRoleLogic) AddSysRole(req *types.AddSysRoleReq) error {
	_, err := l.svcCtx.SysRoleModel.FindOneByUniqueKey(l.ctx, req.UniqueKey)
	if err == model.ErrNotFound {
		if req.ParentId != globalkey.SysTopParentId {
			_, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, req.ParentId)
			if err != nil {
				return errorx2.NewDefaultError(errorx2.ParentRoleIdErrorCode)
			}
		}

		var sysRole = new(model.SysRole)
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
		_, err = l.svcCtx.SysRoleModel.Insert(l.ctx, sysRole)
		if err != nil {
			return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}

		return nil
	} else {

		return errorx2.NewDefaultError(errorx2.AddRoleErrorCode)
	}
}
