package role

import (
	"context"
	"encoding/json"
	errorx2 "zero-zone/pkg/errorx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type GetSysRoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSysRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSysRoleListLogic {
	return &GetSysRoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSysRoleListLogic) GetSysRoleList() (resp *types.SysRoleListResp, err error) {
	sysRoleList, err := l.svcCtx.SysRoleModel.FindAll(l.ctx)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var role types.Role
	roleList := make([]types.Role, 0)
	for _, v := range sysRoleList {
		err := copier.Copy(&role, &v)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		var permMenuIds []int64
		err = json.Unmarshal([]byte(v.PermMenuIds), &permMenuIds)
		role.PermMenuIds = permMenuIds
		roleList = append(roleList, role)
	}

	return &types.SysRoleListResp{
		List: roleList,
	}, nil
}
