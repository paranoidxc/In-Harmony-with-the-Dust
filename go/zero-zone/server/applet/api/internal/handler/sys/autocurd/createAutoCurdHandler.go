package autocurd

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"zero-zone/applet/api/internal/logic/sys/autocurd"
	"zero-zone/applet/api/internal/types"
	"zero-zone/pkg/myvalid"
	"zero-zone/pkg/response"

	"zero-zone/applet/api/internal/svc"
)

// 新增
func CreateAutoCurdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AutoCurdCreateReq
		if err := httpx.Parse(r, &req); err != nil {
		}

		validateErr := myvalid.Validate(&req)
		if validateErr != nil {
			response.Response(w, nil, validateErr)
			return
		}

		l := autocurd.NewCreateAutoCurdLogic(r.Context(), svcCtx)
		err := l.CreateAutoCurd(&req)
		if err != nil {
			response.Response(w, nil, err)
			return
		}

		response.Response(w, nil, nil)
	}
}

// 新增gorm方式生成
func CreateAutoCurdHandlerGorm(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := autocurd.NewCreateAutoCurdLogicGorm(r.Context(), svcCtx)
		err := l.CreateAutoCurd()
		if err != nil {
			response.Response(w, nil, err)
			return
		}

		response.Response(w, nil, nil)
	}
}
