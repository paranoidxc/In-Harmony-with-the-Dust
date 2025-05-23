package demo_curd

import (
	"net/http"
	"zero-zone/pkg/myvalid"
	"zero-zone/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zero-zone/applet/api/internal/logic/feat/demo_curd"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

func DemoCurdUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DemoCurdUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(w, nil, err)
			return
		}

		validateErr := myvalid.Validate(&req)
		if validateErr != nil {
			response.Response(w, nil, validateErr)
			return
		}

		l := demo_curd.NewDemoCurdUpdateLogic(r.Context(), svcCtx)
		err := l.DemoCurdUpdate(&req)
		if err != nil {
			response.Response(w, nil, err)
			return
		}

		response.Response(w, nil, nil)
	}
}
