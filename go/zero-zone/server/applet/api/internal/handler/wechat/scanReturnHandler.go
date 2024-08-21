package wechat

import (
	"net/http"
	"zero-zone/pkg/myvalid"
	"zero-zone/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zero-zone/applet/api/internal/logic/wechat"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

func ScanReturnHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WechatScanReturnReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(w, nil, err)
			return
		}

		validateErr := myvalid.Validate(&req)
		if validateErr != nil {
			response.Response(w, nil, validateErr)
			return
		}

		l := wechat.NewScanReturnLogic(r.Context(), svcCtx)
		resp, err := l.ScanReturn(&req)
		response.ResponseHtml(r.Context(), w, resp.Msg, err)
		/*
				if err != nil {
					response.Response(w, nil, err)
					return
				}
			fmt.Println(resp)
			response.ResponseHtml(r.Context(), w, "ACB")
		*/
		//response.Response(w, resp, nil)
	}
}
