package wechat

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-zone/pkg/utils"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginQRCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

const cacheWechatLoginTokenPrefix = "cache:zeroZone:wechatLoginTokenPrefix:"

func getCacheWechatLoginTokenRedisFullKey(cacheKey string) string {
	return fmt.Sprintf("%s%v", cacheWechatLoginTokenPrefix, cacheKey)
}

func NewLoginQRCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginQRCodeLogic {
	return &LoginQRCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginQRCodeLogic) LoginQRCode(req *types.WechatLoginQRCodeReq) (resp *types.WechatLoginQRCodeResp, err error) {
	token := utils.RandomString(20)
	device := utils.RandomString(10)

	//scheme := "http://abc.com"
	//if l.ctx.c.Request.TLS != nil {
	//	scheme = "https://"
	//absoluteReturnURL := scheme + c.Request.Host + "/wechat/loginRet?token=" + token + "&device=" + device
	absoluteReturnURL := req.Host + "/wechat/loginRet?token=" + token + "&device=" + device
	url := "http://demo.demowechat.com/wxapi?c=login&scope=snsapi_userinfo&url=" + absoluteReturnURL
	logc.Infow(l.ctx, "wechat", logx.Field("url", url))
	qrCode, err := utils.QrcodeWeChat(url)
	if err != nil {
		err = errors.New("获取登录二维码失败")
		return
	}

	cacheKey := getCacheWechatLoginTokenRedisFullKey(token)
	_ = l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, "", l.svcCtx.Config.ThirdPartAllowLoginTokenExpire)

	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCode)
	resp = &types.WechatLoginQRCodeResp{
		Token:         token,
		QRCodeContent: "data:image/png;base64," + qrCodeBase64,
	}

	return
}
