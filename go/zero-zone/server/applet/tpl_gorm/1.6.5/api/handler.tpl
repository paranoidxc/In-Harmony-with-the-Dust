package {{.PkgName}}

import (
	"net/http"
	"zero-zone/pkg/response"
	"zero-zone/pkg/myvalid"


	"github.com/zeromicro/go-zero/rest/httpx"
	{{.ImportPackages}}
)

{{if .HasDoc}}{{.Doc}}{{end}}
func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := httpx.Parse(r, &req); err != nil {
			//response.Response(w, nil, err)
			//return
		}

		validateErr := myvalid.Validate(&req)
  		if validateErr != nil {
	        response.Response(w, nil, validateErr)
            return
        }

		{{end}}l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
		if err != nil {
			response.Response(w, nil, err)
			return
		}

		{{if .HasResp}}response.Response(w, resp, nil){{else}}response.Response(w, nil, nil){{end}}
	}
}
