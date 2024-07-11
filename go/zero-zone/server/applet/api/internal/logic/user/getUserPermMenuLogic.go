package user

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"
	"zero-zone/pkg/utils"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
	"zero-zone/applet/model"
)

type GetUserPermMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserPermMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPermMenuLogic {
	return &GetUserPermMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserPermMenuLogic) GetUserPermMenu() (resp *types.UserPermMenuResp, err error) {
	userId := utils.GetUserId(l.ctx)

	online, err := l.svcCtx.Redis.Get(globalkey.SysOnlineUserCachePrefix + strconv.FormatInt(userId, 10))
	if err != nil || online == "" {
		return nil, errorx2.NewDefaultError(errorx2.AuthErrorCode)
	}

	// 查询用户信息
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var roles []int64
	// 用户所属角色
	err = json.Unmarshal([]byte(user.RoleIds), &roles)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var permMenu []int64
	var userPermMenu []*model.SysPermMenu

	userPermMenu, permMenu, err = l.countUserPermMenu(roles, permMenu)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var menu types.Menu
	menuList := make([]types.Menu, 0)
	newMenuList := make([]types.NewMenu, 0)

	permList := make([]string, 0)
	_, err = l.svcCtx.Redis.Del(globalkey.SysPermMenuCachePrefix + strconv.FormatInt(userId, 10))

	mapParentId2Children := make(map[int64][]int64)

	mapId2Menu := make(map[int64]types.NewMenu)

	for _, v := range userPermMenu {
		if v.IsShow == 0 {
			continue
		}
		err := copier.Copy(&menu, &v)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}

		tmp := []types.NewMenu{}
		isShow := true
		if v.IsShow == 0 {
			isShow = false
		}

		keepAlive := false
		if v.KeepAlive == 1 {
			keepAlive = true
		}

		newMenu := types.NewMenu{
			Title:            menu.Name,
			MenuId:           menu.Id,
			NewTest:          strings.Replace(menu.ViewPath, "views/", "", -1), // vue页面地址
			Component:        strings.Replace(menu.ViewPath, "views/", "", -1), // vue页面地址
			Path:             menu.Router,
			Icon:             menu.Icon,
			IsHasPerm:        true,
			IsShow:           isShow,
			KeepAlive:        keepAlive,
			IsShowBreadcrumb: true,
			ParentId:         0,
			ParentName:       "ParentName",
			OrderNum:         v.OrderNum,
			Children:         tmp,
		}
		mapId2Menu[menu.Id] = newMenu
		if menu.Type != globalkey.SysDefaultPermType {
			menuList = append(menuList, menu)
			if menu.ParentId == 0 {
				newMenuList = append(newMenuList, newMenu)
				/*
					newMenuList = append(newMenuList, types.NewMenu{
						Title:            menu.Name,
						MenuId:           menu.Id,
						NewTest:          strings.Replace(menu.ViewPath, "views/", "", -1), // vue页面地址
						Component:        strings.Replace(menu.ViewPath, "views/", "", -1), // vue页面地址
						Path:             menu.Router,
						Icon:             "Odometer",
						IsHasPerm:        true,
						IsShow:           true,
						IsShowBreadcrumb: true,
						ParentId:         0,
						ParentName:       "ParentName",
						Children:         tmp,
					})
				*/
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

		var permArray []string
		err = json.Unmarshal([]byte(v.Perms), &permArray)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}

		for _, p := range permArray {
			_, err := l.svcCtx.Redis.Sadd(globalkey.SysPermMenuCachePrefix+strconv.FormatInt(userId, 10), globalkey.SysPermMenuPrefix+p)
			if err != nil {
				return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
			}
			permList = append(permList, "/"+p)
		}

	}
	//fmt.Println("len:", len(newMenuList))
	//fmt.Println("%+v:", mapParentId2Children)
	//fmt.Println("%+v:", mapId2Menu)
	for idx, tmpMenu := range newMenuList {
		if tmpMenu.IsShow {
			tmpList := []types.NewMenu{}
			val, ok := mapParentId2Children[tmpMenu.MenuId]
			if ok {
				for _, menuId := range val {
					tmpMenu := mapId2Menu[menuId]
					if tmpMenu.IsShow {
						tmpList = append(tmpList, tmpMenu)
					}
				}
			}

			tmpMenu.Children = tmpList
			newMenuList[idx] = tmpMenu
		}
	}

	sort.Slice(newMenuList, func(i, j int) bool {
		return newMenuList[i].OrderNum > newMenuList[j].OrderNum
	})

	return &types.UserPermMenuResp{PermissionTreeList: newMenuList, Menus: menuList, Perms: utils.ArrayUniqueValue[string](permList)}, nil
}

func (l *GetUserPermMenuLogic) countUserPermMenu(roles []int64, permMenu []int64) ([]*model.SysPermMenu, []int64, error) {
	if utils.ArrayContainValue(roles, globalkey.SysSuperRoleId) {
		sysPermMenus, err := l.svcCtx.SysPermMenuModel.FindAll(l.ctx)
		if err != nil {
			return nil, permMenu, err
		}

		return sysPermMenus, permMenu, nil
	} else {
		for _, roleId := range roles {
			// 查询角色信息
			role, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, roleId)
			if err != nil && err != model.ErrNotFound {
				return nil, permMenu, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
			}

			var perms []int64
			// 角色所拥有的权限id
			err = json.Unmarshal([]byte(role.PermMenuIdsAll), &perms)
			if err != nil {
				return nil, permMenu, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
			}

			if role.Status != 0 {
				permMenu = append(permMenu, perms...)
			}
			// 汇总用户所属角色权限id
			permMenu = l.getRolePermMenu(permMenu, roleId)
		}

		// 过滤重复的权限id
		permMenu = utils.ArrayUniqueValue[int64](permMenu)
		var ids string
		for i, id := range permMenu {
			if i == 0 {
				ids = strconv.FormatInt(id, 10)
				continue
			}
			ids = ids + "," + strconv.FormatInt(id, 10)
		}

		if len(ids) == 0 {
			return nil, permMenu, nil
		}

		// 根据权限id获取具体权限
		sysPermMenus, err := l.svcCtx.SysPermMenuModel.FindByIds(l.ctx, ids)
		if err != nil {
			return nil, permMenu, err
		}

		return sysPermMenus, permMenu, nil
	}
}

func (l *GetUserPermMenuLogic) getRolePermMenu(perms []int64, roleId int64) []int64 {
	roles, err := l.svcCtx.SysRoleModel.FindSubRole(l.ctx, roleId)
	if err != nil && err != model.ErrNotFound {
		return perms
	}

	for _, role := range roles {
		var subPerms []int64
		err = json.Unmarshal([]byte(role.PermMenuIds), &subPerms)
		perms = append(perms, subPerms...)
		perms = l.getRolePermMenu(perms, role.Id)
	}

	return perms
}
