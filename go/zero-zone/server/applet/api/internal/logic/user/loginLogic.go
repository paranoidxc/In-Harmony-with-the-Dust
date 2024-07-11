package user

import (
	"context"
	"net/http"
	"strconv"
	"time"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"
	"zero-zone/pkg/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
	"zero-zone/applet/model"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq, r *http.Request) (resp *types.LoginResp, err error) {
	/*
		verifyCode, _ := l.svcCtx.Redis.Get(globalkey.SysLoginCaptchaCachePrefix + req.CaptchaId)
		if verifyCode != req.VerifyCode {
			return nil, errorx.NewDefaultError(errorx.CaptchaErrorCode)
		}
	*/

	sysUser, err := l.svcCtx.SysUserModel.FindOneByAccount(l.ctx, req.Username)
	if err != nil {
		//return nil, errorx2.NewDefaultError(errorx2.AccountErrorCode)
		return nil, errorx2.NewDefaultError(errorx2.LoginErrorCode)
	}

	if sysUser.Password != utils.MD5(req.Password+l.svcCtx.Config.Salt) {
		//return nil, errorx2.NewDefaultError(errorx2.PasswordErrorCode)
		return nil, errorx2.NewDefaultError(errorx2.LoginErrorCode)
		//	return nil, errors.New("账户名称或密码不正确。")
	}

	if sysUser.Status != globalkey.SysEnable {
		return nil, errorx2.NewDefaultError(errorx2.AccountDisableErrorCode)
	}

	if sysUser.Id != globalkey.SysSuperUserId {
		/*
			dept, _ := l.svcCtx.SysDeptModel.FindOne(l.ctx, sysUser.DeptId)
			if dept.Status == globalkey.SysDisable {
				return nil, errorx2.NewDefaultError(errorx2.AccountDisableErrorCode)
			}
		*/
	}

	token, _ := l.getJwtToken(sysUser.Id)
	_, err = l.svcCtx.Redis.Del(req.CaptchaId)

	loginLog := model.SysLog{
		UserId: sysUser.Id,
		Ip:     r.Header.Get("X-Forwarded-For"),
		Uri:    r.RequestURI,
		Type:   1,
		Status: 1,
	}
	_, err = l.svcCtx.SysLogModel.Insert(l.ctx, &loginLog)

	err = l.svcCtx.Redis.Setex(globalkey.SysOnlineUserCachePrefix+strconv.FormatInt(sysUser.Id, 10), "1", int(l.svcCtx.Config.JwtAuth.AccessExpire))
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return &types.LoginResp{
		Token:      token,
		TokenName:  "Authorization",
		TokenValue: token,
	}, nil
}

func (l *LoginLogic) getJwtToken(userId int64) (string, error) {
	iat := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + l.svcCtx.Config.JwtAuth.AccessExpire
	claims["iat"] = iat
	claims[globalkey.SysJwtUserId] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(l.svcCtx.Config.JwtAuth.AccessSecret))
}
