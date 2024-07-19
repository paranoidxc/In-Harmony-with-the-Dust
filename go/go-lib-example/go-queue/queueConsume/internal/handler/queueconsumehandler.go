package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"queueConsume/internal/logic"
	"queueConsume/internal/svc"
	"queueConsume/internal/types"
)

func QueueConsumeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewQueueConsumeLogic(r.Context(), svcCtx)
		resp, err := l.QueueConsume(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
