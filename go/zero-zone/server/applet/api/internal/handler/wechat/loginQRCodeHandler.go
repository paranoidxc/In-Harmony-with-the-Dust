package wechat

import (
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"zero-zone/applet/api/internal/types"
	"zero-zone/pkg/response"
	"zero-zone/pkg/utils"

	"zero-zone/applet/api/internal/logic/wechat"
	"zero-zone/applet/api/internal/svc"
)

func LoginQRCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := wechat.NewLoginQRCodeLogic(r.Context(), svcCtx)

		host := utils.Host(r)
		logc.Infow(r.Context(), "host", logx.Field("host", host))
		req := &types.WechatLoginQRCodeReq{
			Host: host,
		}
		resp, err := l.LoginQRCode(req)
		if err != nil {
			response.Response(w, nil, err)
			return
		}

		response.Response(w, resp, nil)
	}
}
