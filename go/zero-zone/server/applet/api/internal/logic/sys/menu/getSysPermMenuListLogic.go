package menu

import (
	"context"
	"encoding/json"
	"strconv"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type GetSysPermMenuListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSysPermMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSysPermMenuListLogic {
	return &GetSysPermMenuListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSysPermMenuListLogic) GetSysPermMenuList() (resp *types.SysPermMenuListResp, err error) {
	permMenus, err := l.svcCtx.SysPermMenuModel.FindAll(l.ctx)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var menu types.PermMenu
	PermMenuList := make([]types.PermMenu, 0)

	mapId2Menu := make(map[int64]types.PermMenu)
	mapParentId2Children := make(map[int64][]int64)

	for _, v := range permMenus {
		err := copier.Copy(&menu, &v)
		menu.Label = menu.Name
		menu.Value = menu.Id
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}

		var perms []string
		err = json.Unmarshal([]byte(v.Perms), &perms)
		menu.Perms = perms

		mapId2Menu[menu.Id] = menu

		if menu.ParentId == 0 {
			PermMenuList = append(PermMenuList, menu)
		} else {
			val, ok := mapParentId2Children[menu.ParentId]
			if ok {
				val := append(val, menu.Id)
				mapParentId2Children[menu.ParentId] = val
			} else {
				mapParentId2Children[menu.ParentId] = []int64{menu.Id}
			}
		}
	}
	/*
		for idx, tmpMenu := range PermMenuList {
			tmpList := []types.PermMenu{}
			val, ok := mapParentId2Children[tmpMenu.Id]
			if ok {
				for _, menuId := range val {
					tmpMenu := mapId2Menu[menuId]
					tmpList = append(tmpList, tmpMenu)
				}
			}

			tmpMenu.Children = tmpList
			PermMenuList[idx] = tmpMenu
		}
	*/

	PermMenuList = buildTreeList(PermMenuList, mapParentId2Children, mapId2Menu)

	return &types.SysPermMenuListResp{List: PermMenuList}, nil
}

func buildTreeList(list []types.PermMenu,
	mapParentId2Children map[int64][]int64,
	mapId2Menu map[int64]types.PermMenu) []types.PermMenu {

	for idx, menu := range list {
		childrenList := []types.PermMenu{}
		val, ok := mapParentId2Children[menu.Id]
		if ok {
			for _, menuId := range val {
				tmpMenu := mapId2Menu[menuId]
				childrenList = append(childrenList, tmpMenu)
			}
		}

		if len(childrenList) > 0 {
			childrenList = buildTreeList(childrenList, mapParentId2Children, mapId2Menu)
		}

		menu.Children = childrenList
		list[idx] = menu
	}

	return list
}

func (l *GetSysPermMenuListLogic) getCurrentUserPermMenuIds(currentUserId int64) (ids []int64) {
	var currentPermMenuIds []int64
	if currentUserId != globalkey.SysSuperUserId {
		var currentUserRoleIds []int64
		var roleIds []int64
		currentUser, _ := l.svcCtx.SysUserModel.FindOne(l.ctx, currentUserId)
		_ = json.Unmarshal([]byte(currentUser.RoleIds), &currentUserRoleIds)
		roleIds = append(roleIds, currentUserRoleIds...)
		var ids string
		for i, v := range roleIds {
			if i == 0 {
				ids = strconv.FormatInt(v, 10)
			}
			ids = ids + "," + strconv.FormatInt(v, 10)
		}

		sysRoles, _ := l.svcCtx.SysRoleModel.FindByIds(l.ctx, ids)
		var rolePermMenus []int64
		for _, v := range sysRoles {
			err := json.Unmarshal([]byte(v.PermMenuIds), &rolePermMenus)
			if err != nil {
				return nil
			}
			currentPermMenuIds = append(currentPermMenuIds, rolePermMenus...)
		}
	}

	return currentPermMenuIds
}
