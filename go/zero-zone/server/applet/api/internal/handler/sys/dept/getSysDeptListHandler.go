package dept

import (
	"net/http"
	"zero-zone/applet/api/internal/types"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zero-zone/applet/api/internal/logic/sys/dept"
	"zero-zone/applet/api/internal/svc"
)

func GetSysDeptListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SysDeptPageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx2.NewHandlerError(errorx2.ParamErrorCode, err.Error()))
			return
		}
		l := dept.NewGetSysDeptListLogic(r.Context(), svcCtx)
		//resp, err := l.GetSysDeptList()
		resp, err := l.GetSysDeptPage(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		response.Response(w, resp, err)
	}
}
