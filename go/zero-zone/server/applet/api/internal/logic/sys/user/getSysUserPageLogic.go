package user

import (
	"context"
	"strconv"
	"strings"
	errorx2 "zero-zone/pkg/errorx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
	"zero-zone/applet/model"
)

type GetSysUserPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSysUserPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSysUserPageLogic {
	return &GetSysUserPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSysUserPageLogic) GetSysUserPage(req *types.SysUserPageReq) (resp *types.SysUserPageResp, err error) {
	s := strconv.FormatInt(req.DeptId, 10)
	deptIds := l.getDeptIds(s, req.DeptId)

	users, err := l.svcCtx.SysUserModel.FindPage(l.ctx, req.Page, req.Limit, deptIds)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var user types.User
	//var userProfession types.UserProfession
	//var userJob types.UserJob
	//var userDept types.UserDept
	userList := make([]types.User, 0)
	for _, v := range users {
		err := copier.Copy(&user, &v)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}

		//userProfession.Id = v.ProfessionId
		//userProfession.Name = v.Profession

		//userJob.Id = v.JobId
		//userJob.Name = v.Job

		//userDept.Id = v.DeptId
		//userDept.Name = v.Dept

		var userRole types.UserRole
		var roles []types.UserRole
		var roleNameArr []string
		var roleIdArr []string
		roleNameArr = strings.Split(v.Roles, ",")
		roleIdArr = strings.Split(v.RoleIds, ",")
		for i, n := range roleNameArr {
			userRole.Name = n
			userRole.Id, _ = strconv.ParseInt(roleIdArr[i], 10, 64)
			roles = append(roles, userRole)
		}

		//user.Profession = userProfession
		//user.Job = userJob
		//user.Dept = userDept
		user.Roles = roles

		userList = append(userList, user)
	}

	total, err := l.svcCtx.SysUserModel.FindCountByDeptIds(l.ctx, deptIds)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	pagination := types.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	return &types.SysUserPageResp{
		List:       userList,
		Pagination: pagination,
	}, nil
}

func (l *GetSysUserPageLogic) getDeptIds(deptId string, id int64) string {
	deptList, err := l.svcCtx.SysDeptModel.FindSubDept(l.ctx, id)
	if err != nil && err != model.ErrNotFound {
		return deptId
	}

	for _, v := range deptList {
		deptId = deptId + "," + strconv.FormatInt(v.Id, 10)
		deptId = l.getDeptIds(deptId, v.Id)
	}

	return deptId
}
