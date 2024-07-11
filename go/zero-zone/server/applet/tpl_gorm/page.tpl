
func (l *{{ .Name }}PageLogic) {{ .Name }}Page(req *types.{{ .Name }}PageReq) (resp *types.{{ .Name }}PageResp, err error) {
    where := " 1 "
    /*
    {{- range $i, $v := .VueFields }}
    if len(strings.TrimSpace(req.{{ $v.Name }})) > 0 {
        where = where + fmt.Sprintf(" AND {{ $v.Column }} LIKE '%s'", "%"+strings.TrimSpace(req.{{ $v.Name }})+"%")
    }
    {{- end }}
    */

    feat{{ .Name }}Page, err := l.svcCtx.Feat{{ .Name }}Model.FindPageByWhere(where, req.Page, req.Limit)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var item types.{{ .Name }}
	{{ .Name }}Page := make([]types.{{ .Name }}, 0)
	for _, v := range feat{{ .Name }}Page {
		err := copier.Copy(&item, &v)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		{{ .Name }}Page = append({{ .Name }}Page, item)
	}

	total, err := l.svcCtx.Feat{{ .Name }}Model.FindPageByWhereCount(where)
    if err != nil {
         return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
    }

    pagination := types.Pagination{
         Page:  req.Page,
         Limit: req.Limit,
         Total: total,
    }

	return &types.{{ .Name }}PageResp{
		List: {{ .Name }}Page,
   		Pagination: pagination,
	}, nil
}