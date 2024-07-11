package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type LogActionMiddleware struct {
}

func NewLogActionMiddleware() *LogActionMiddleware {
	return &LogActionMiddleware{}
}

func (m *LogActionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body []byte
		if r.Method == "GET" {
			query := r.URL.RawQuery
			query, _ = url.QueryUnescape(query)
			split := strings.Split(query, "&")
			m := make(map[string]string)
			for _, v := range split {
				kv := strings.Split(v, "=")
				if len(kv) == 2 {
					m[kv[0]] = kv[1]
				}
			}
			body, _ = json.Marshal(&m)
		} else {
			var err error
			body, err = io.ReadAll(r.Body)
			if err != nil {
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		}
		logc.Infow(r.Context(),
			"==========Req Info==========",
			logx.Field("RequestURI", r.RequestURI),
			logx.Field("body", string(body)),
		)
		respWrite := responseBodyWriter{w, bytes.NewBuffer(body)}
		next(respWrite, r)
		logc.Infow(r.Context(),
			"==========Resp Info==========",
			logx.Field("body", respWrite.body.String()),
		)
	}
}

type responseBodyWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
