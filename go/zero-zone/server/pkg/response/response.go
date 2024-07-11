package response

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"net/http"
	"net/url"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	TraceId string      `json:"traceId"`
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
}

func SendFile(ctx context.Context, w http.ResponseWriter, downloadName string, b []byte) {
	contentDisposition := fmt.Sprintf("attachment; filename=%s", url.QueryEscape(downloadName))
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	w.Header().Set("Content-Disposition", contentDisposition)
	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	w.Write(b)
}

func ResponseWithCtx(ctx context.Context, w http.ResponseWriter, resp interface{}, err error) {
	tracer := otel.GetTracerProvider().Tracer(trace.TraceName)
	_, span := tracer.Start(ctx, "ResponseWithCtx")
	var body Body
	body.TraceId = span.SpanContext().TraceID().String()
	if err != nil {
		body.Code = 0
		msg := err.Error()
		if strings.Contains(msg, "mismatch") {
			body.Msg = "格式错误:" + err.Error()
		} else {
			body.Msg = err.Error()
		}
	} else {
		body.Code = 200
		body.Msg = "success"
		body.Data = resp
	}
	httpx.OkJson(w, body)
}

func Response(w http.ResponseWriter, resp interface{}, err error) {
	var body Body
	if err != nil {
		body.Code = 0
		msg := err.Error()
		if strings.Contains(msg, "mismatch") {
			body.Msg = "格式错误:" + err.Error()
		} else {
			body.Msg = err.Error()
		}
	} else {
		body.Code = 200
		body.Msg = "success"
		body.Data = resp
	}
	httpx.OkJson(w, body)
}
