package sys_region

import (
	"net/http"
	"zero-zone/pkg/response"

	"zero-zone/applet/api/internal/logic/feat/sys_region"
	"zero-zone/applet/api/internal/svc"
)

func SysRegionTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := sys_region.NewSysRegionTreeLogic(r.Context(), svcCtx)
		resp, err := l.SysRegionTree()
		if err != nil {
			response.Response(w, nil, err)
			return
		}

		response.Response(w, resp, nil)
	}
}
