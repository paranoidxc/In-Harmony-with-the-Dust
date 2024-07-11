package middleware

import (
	"net/http"
	"strconv"
	"strings"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"
	"zero-zone/pkg/utils"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type PermMenuAuthMiddleware struct {
	Redis *redis.Redis
}

func NewPermMenuAuthMiddleware(r *redis.Redis) *PermMenuAuthMiddleware {
	return &PermMenuAuthMiddleware{
		Redis: r,
	}
}

func (m *PermMenuAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(r.Header.Get("Authorization")) > 0 {
			userId := utils.GetUserId(r.Context())
			online, err := m.Redis.Get(globalkey.SysOnlineUserCachePrefix + strconv.FormatInt(userId, 10))
			if err != nil || online == "" {
				httpx.Error(w, errorx2.NewDefaultError(errorx2.AuthErrorCode))
				var erring any
				erring = "Auth fail"
				panic(erring)
			}

			uri := strings.Split(r.RequestURI, "?")
			is, err := m.Redis.Sismember(globalkey.SysPermMenuCachePrefix+strconv.FormatInt(userId, 10), uri[0])
			if err != nil || is != true {
				httpx.Error(w, errorx2.NewDefaultError(errorx2.NotPermMenuErrorCode))
			} else {
				next(w, r)
			}
		}
	}
}
