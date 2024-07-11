package dict

import (
	"net/http"
	"zero-zone/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zero-zone/applet/api/internal/logic/config/dict"
	"zero-zone/applet/api/internal/svc"
)

func GetConfigDictListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := dict.NewGetConfigDictListLogic(r.Context(), svcCtx)
		resp, err := l.GetConfigDictList()
		if err != nil {
			httpx.Error(w, err)
			return
		}

		response.Response(w, resp, err)
	}
}
