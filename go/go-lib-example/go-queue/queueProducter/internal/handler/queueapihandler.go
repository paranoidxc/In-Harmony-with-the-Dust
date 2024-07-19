package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"queueapi/internal/logic"
	"queueapi/internal/svc"
	"queueapi/internal/types"
)

func QueueapiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewQueueapiLogic(r.Context(), svcCtx)
		resp, err := l.Queueapi(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
