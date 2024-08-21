package wechat

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logc"
	"strconv"
	"time"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
	"zero-zone/applet/model"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginCheckLogic {
	return &LoginCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginCheckLogic) LoginCheck(req *types.WechatLoginCheckReq) (resp *types.WechatLoginCheckResp, err error) {
	openID, err := l.svcCtx.Redis.GetCtx(l.ctx, getCacheWechatLoginTokenRedisFullKey(req.Token))
	if err != nil || openID == "" {
		err = errors.New("未登录")
		return
	}
	resp = &types.WechatLoginCheckResp{}
	// 根据 openID 取得用户的信息 判断是否已经绑定
	SysUserExt, xerr := l.svcCtx.SysUserModel.FindOneByWechatOpenID(l.ctx, openID)
	if xerr == nil {
		logc.Infow(l.ctx, "已绑定, 生成token信息, 显示信息让用户回到页面")
		resp.Status = 1

		// login logic
		token, _ := l.getJwtToken(SysUserExt.Id)
		err = l.svcCtx.Redis.Setex(globalkey.SysOnlineUserCachePrefix+strconv.FormatInt(SysUserExt.Id, 10), "1", int(l.svcCtx.Config.JwtAuth.AccessExpire))
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}

		resp.Token = token
		resp.TokenName = "Authorization"
		resp.TokenValue = token
	} else if xerr == model.ErrNotFound {
		logc.Infow(l.ctx, "未绑定, 显示信息让用户回到页面进行账号绑定")
		resp.Status = 2
	} else {
		logc.Errorw(l.ctx, "系统错误", logc.Field("err", xerr))
		return nil, errors.New("系统错误:" + xerr.Error())
	}

	return
}

func (l *LoginCheckLogic) getJwtToken(userId int64) (string, error) {
	iat := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + l.svcCtx.Config.JwtAuth.AccessExpire
	claims["iat"] = iat
	claims[globalkey.SysJwtUserId] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(l.svcCtx.Config.JwtAuth.AccessSecret))
}
