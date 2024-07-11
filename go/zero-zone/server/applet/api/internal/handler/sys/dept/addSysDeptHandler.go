package dept

import (
	"fmt"
	"net/http"
	"zero-zone/pkg/myvalid"
	"zero-zone/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zero-zone/applet/api/internal/logic/sys/dept"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

func AddSysDeptHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddSysDeptReq
		if err := httpx.Parse(r, &req); err != nil {
			//httpx.Error(w, errorx2.NewHandlerError(errorx2.ParamErrorCode, err.Error()))
			//return
		}
		validateErr := myvalid.Validate(&req)
		if validateErr != nil {
			fmt.Println("validEr", validateErr)
			response.Response(w, nil, validateErr)
			return
		}

		/*
			validate := validator.New()
			validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
				name := fld.Tag.Get("label")
				return name
			})

			trans, _ := ut.New(zh.New()).GetTranslator("zh")
			validateErr := translations.RegisterDefaultTranslations(validate, trans)
			if validateErr = validate.StructCtx(r.Context(), req); validateErr != nil {
				for _, err := range validateErr.(validator.ValidationErrors) {
					httpx.Error(w, errorx2.NewHandlerError(errorx2.ParamErrorCode, errors.New(err.Translate(trans)).Error()))
					return
				}
			}

		*/

		l := dept.NewAddSysDeptLogic(r.Context(), svcCtx)
		err := l.AddSysDept(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		response.Response(w, nil, err)
	}
}
