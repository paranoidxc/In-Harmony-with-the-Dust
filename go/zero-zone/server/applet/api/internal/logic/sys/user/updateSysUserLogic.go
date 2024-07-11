package user

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"strconv"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"

	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type UpdateSysUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysUserLogic {
	return &UpdateSysUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysUserLogic) UpdateSysUser(req *types.UpdateSysUserReq) error {
	/*
		_, err := l.svcCtx.SysDeptModel.FindOne(l.ctx, req.DeptId)
		if err != nil {
			return errorx2.NewDefaultError(errorx2.DeptIdErrorCode)
		}

		_, err = l.svcCtx.SysProfessionModel.FindOne(l.ctx, req.ProfessionId)
		if err != nil {
			return errorx2.NewDefaultError(errorx2.ProfessionIdErrorCode)
		}

		_, err = l.svcCtx.SysJobModel.FindOne(l.ctx, req.JobId)
		if err != nil {
			return errorx2.NewDefaultError(errorx2.JobIdErrorCode)
		}
	*/

	editUser, err := l.svcCtx.SysUserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errorx2.NewDefaultError(errorx2.UserIdErrorCode)
	}

	err = copier.Copy(editUser, req)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	bytes, err := json.Marshal(req.RoleIds)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}
	editUser.RoleIds = string(bytes)

	err = l.svcCtx.SysUserModel.Update(l.ctx, editUser)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	_, err = l.svcCtx.Redis.Del(globalkey.SysPermMenuCachePrefix + strconv.FormatInt(editUser.Id, 10))
	_, err = l.svcCtx.Redis.Del(globalkey.SysOnlineUserCachePrefix + strconv.FormatInt(editUser.Id, 10))
	//fmt.Println("okay")
	return nil
}
