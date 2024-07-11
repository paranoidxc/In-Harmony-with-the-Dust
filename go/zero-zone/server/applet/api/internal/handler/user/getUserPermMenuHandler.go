package user

import (
	"net/http"
	"zero-zone/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zero-zone/applet/api/internal/logic/user"
	"zero-zone/applet/api/internal/svc"
)

func GetUserPermMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetUserPermMenuLogic(r.Context(), svcCtx)
		resp, err := l.GetUserPermMenu()
		if err != nil {
			httpx.Error(w, err)
			return
		}

		response.Response(w, resp, err)
	}
}
