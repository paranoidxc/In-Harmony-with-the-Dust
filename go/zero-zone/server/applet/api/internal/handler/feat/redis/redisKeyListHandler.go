package redis

import (
	"net/http"
	"zero-zone/pkg/response"

	"zero-zone/applet/api/internal/logic/feat/redis"
	"zero-zone/applet/api/internal/svc"
)

func RedisKeyListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := redis.NewRedisKeyListLogic(r.Context(), svcCtx)
		resp, err := l.RedisKeyList()
		if err != nil {
			response.Response(w, nil, err)
			return
		}

		response.Response(w, resp, nil)
	}
}
