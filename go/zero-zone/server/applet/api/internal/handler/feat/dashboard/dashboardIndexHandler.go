package dashboard

import (
	"net/http"
	"zero-zone/pkg/myvalid"
	"zero-zone/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zero-zone/applet/api/internal/logic/feat/dashboard"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

func DashboardIndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DashboardIndexReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(w, nil, err)
			return
		}

		validateErr := myvalid.Validate(&req)
		if validateErr != nil {
			response.Response(w, nil, validateErr)
			return
		}

		l := dashboard.NewDashboardIndexLogic(r.Context(), svcCtx)
		resp, err := l.DashboardIndex(&req)
		if err != nil {
			response.Response(w, nil, err)
			return
		}

		response.Response(w, resp, nil)
	}
}
