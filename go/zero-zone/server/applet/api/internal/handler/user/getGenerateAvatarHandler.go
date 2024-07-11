package user

import (
	"net/http"
	"zero-zone/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zero-zone/applet/api/internal/logic/user"
	"zero-zone/applet/api/internal/svc"
)

func GetGenerateAvatarHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetGenerateAvatarLogic(r.Context(), svcCtx)
		resp, err := l.GetGenerateAvatar()
		if err != nil {
			httpx.Error(w, err)
			return
		}

		response.Response(w, resp, err)
	}
}
