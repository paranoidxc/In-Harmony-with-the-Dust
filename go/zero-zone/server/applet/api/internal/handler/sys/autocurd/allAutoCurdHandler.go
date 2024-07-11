package autocurd

import (
	"net/http"
	"zero-zone/pkg/response"

	"zero-zone/applet/api/internal/logic/sys/autocurd"
	"zero-zone/applet/api/internal/svc"
)

func AllAutoCurdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := autocurd.NewAllAutoCurdLogic(r.Context(), svcCtx)
		resp, err := l.AllAutoCurd()
		if err != nil {
			response.Response(w, nil, err)
			return
		}

		response.Response(w, resp, nil)
	}
}
