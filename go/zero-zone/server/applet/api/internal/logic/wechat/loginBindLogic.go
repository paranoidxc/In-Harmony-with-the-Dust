package wechat

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
	"zero-zone/applet/model"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginBindLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginBindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginBindLogic {
	return &LoginBindLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginBindLogic) LoginBind(req *types.WechatLoginBindReq) (resp *types.WechatLoginBindResp, err error) {
	openID, err := l.svcCtx.Redis.GetCtx(l.ctx, getCacheWechatLoginTokenRedisFullKey(req.Token))
	if err != nil || openID == "" {
		err = errors.New("未登录")
		return
	}

	// 根据 openID 取得用户的信息 判断是否已经绑定
	sysUserExt, xerr := l.svcCtx.SysUserModel.FindOneByWechatOpenID(l.ctx, openID)
	if xerr == nil {
		logc.Infow(l.ctx, "该微信已绑定,请接绑后再操作")
		return nil, errors.New("该微信已绑定,请接绑后再操作")
	} else if xerr == model.ErrNotFound {
		//logc.Infow(l.ctx, "未绑定, 显示信息让用户回到页面进行账号绑定")
	} else {
		logc.Errorw(l.ctx, "系统错误", logc.Field("err", xerr))
		return nil, errors.New("系统错误:" + xerr.Error())
	}

	sysUser, err := l.svcCtx.SysUserModel.FindOneByAccount(l.ctx, req.Username)
	if err != nil {
		return nil, errorx2.NewDefaultError(errorx2.LoginErrorCode)
	}

	logc.Infow(l.ctx, "用户", logc.Field("user", sysUser))
	if sysUser.Password != utils.MD5(req.Password+l.svcCtx.Config.Salt) {
		return nil, errorx2.NewDefaultError(errorx2.LoginErrorCode)
	}

	sysUserExt = &model.SysUserExt{}
	xerr = copier.Copy(sysUserExt, sysUser)
	if xerr != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, xerr.Error())
	}

	sysUserExt.WechatOpenID = openID
	xerr = l.svcCtx.SysUserModel.UpdateExt(l.ctx, sysUserExt)
	if xerr != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, xerr.Error())
	}

	_, xerr = l.svcCtx.Redis.DelCtx(l.ctx, getCacheWechatLoginTokenRedisFullKey(req.Token))
	if xerr != nil {
		logc.Errorw(l.ctx, "删除缓存失败", logc.Field("err", xerr))
	}

	return
}
