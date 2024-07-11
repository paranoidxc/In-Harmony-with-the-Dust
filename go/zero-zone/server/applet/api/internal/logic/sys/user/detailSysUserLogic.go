package user

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"strings"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type DetailSysUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailSysUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailSysUserLogic {
	return &DetailSysUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailSysUserLogic) DetailSysUser(req *types.DetailSysUserReq) (resp *types.DetailSysUserResp, err error) {
	resp = &types.DetailSysUserResp{}
	item, err := l.svcCtx.SysUserModel.FindOneDetail(l.ctx, req.Id)
	err = copier.Copy(resp, item)
	if err != nil {
		logx.Error("复制结果失败", err)
		return nil, err
	}

	var userRole types.UserRole
	var roles []types.UserRole
	var roleNameArr []string
	var roleIdArr []string

	roleNameArr = strings.Split(item.Roles, ",")
	roleIdArr = strings.Split(item.RoleIds, ",")
	for i, n := range roleNameArr {
		userRole.Name = n
		userRole.Id, _ = strconv.ParseInt(roleIdArr[i], 10, 64)
		roles = append(roles, userRole)
	}
	resp.Roles = roles

	return
}
