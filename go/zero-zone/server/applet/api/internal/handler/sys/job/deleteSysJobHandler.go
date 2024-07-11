package job

import (
	"errors"
	"net/http"
	"reflect"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/response"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/zeromicro/go-zero/rest/httpx"
	"zero-zone/applet/api/internal/logic/sys/job"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

func DeleteSysJobHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteSysJobReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx2.NewHandlerError(errorx2.ParamErrorCode, err.Error()))
			return
		}

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

		l := job.NewDeleteSysJobLogic(r.Context(), svcCtx)
		err := l.DeleteSysJob(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		response.Response(w, nil, err)
	}
}
