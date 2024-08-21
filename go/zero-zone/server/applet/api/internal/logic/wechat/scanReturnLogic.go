package wechat

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
	"zero-zone/applet/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScanReturnLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScanReturnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScanReturnLogic {
	return &ScanReturnLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScanReturnLogic) ScanReturn(req *types.WechatScanReturnReq) (resp *types.WechatScanReturnResp, err error) {
	logc.Infow(l.ctx, "扫码登录返回信息", logx.Field("req", req))

	cacheKey := getCacheWechatLoginTokenRedisFullKey(req.Token)
	_ = l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, req.OpenID, l.svcCtx.Config.ThirdPartAllowLoginTokenExpire)

	resp = &types.WechatScanReturnResp{}
	// 根据 openID 取得用户的信息 判断是否已经绑定
	_, xerr := l.svcCtx.SysUserModel.FindOneByWechatOpenID(l.ctx, req.OpenID)
	if xerr == nil {
		resp.Msg = "该微信已绑定账号, 请回到管理系统继续操作"
		logc.Infow(l.ctx, "已绑定, 显示信息让用户回到页面")
	} else if xerr == model.ErrNotFound {
		resp.Msg = "该微信未绑定账号, 请回到管理系统绑定账号"
		logc.Infow(l.ctx, "未绑定, 显示信息让用户回到页面进行账号绑定")
	} else {
		resp.Msg = "非预期错误:" + err.Error()
		logc.Errorw(l.ctx, "非预期错误", logc.Field("err", err))
	}

	return
}
