package user

import (
	"net/http"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/myvalid"
	"zero-zone/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zero-zone/applet/api/internal/logic/sys/user"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

func AddSysUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddSysUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx2.NewHandlerError(errorx2.ParamErrorCode, err.Error()))
			return
		}

		validateErr := myvalid.Validate(&req)
		if validateErr != nil {
			response.Response(w, nil, validateErr)
			return
		}

		l := user.NewAddSysUserLogic(r.Context(), svcCtx)
		err := l.AddSysUser(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		response.Response(w, nil, err)
	}
}
