
func (l *{{ .Name }}PageLogic) {{ .Name }}Page(req *types.{{ .Name }}PageReq) (resp *types.{{ .Name }}PageResp, err error) {
    where := []string{"1"}
    if req != nil {
        if req.IncludeDeleted == 0 {
            where = append(where, "t.deleted_at = '0000-00-00 00:00:00'")
        }
    }
    /*
    {{- range $i, $v := .VueFields }}
    if len(strings.TrimSpace(req.{{ $v.Name }})) > 0 {
        where = append(where, fmt.Sprintf("{{ $v.Column }} LIKE '%s'", "%"+strings.TrimSpace(req.{{ $v.Name }})+"%"))
    }
    {{- end }}
    */

    whereStr := strings.Join(where, " AND ")
    feat{{ .Name }}Page, err := l.svcCtx.Feat{{ .Name }}Model.FindPageByWhere(l.ctx, whereStr, req.Page, req.Limit)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var item types.{{ .Name }}
	{{ .Name }}Page := make([]types.{{ .Name }}, 0)
	for _, v := range feat{{ .Name }}Page {
		err := copier.Copy(&item, &v)
		item.CreatedAt = utils.Time2Str(v.CreatedAt)
       	item.UpdatedAt = utils.Time2Str(v.UpdatedAt)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		{{ .Name }}Page = append({{ .Name }}Page, item)
	}

	total, err := l.svcCtx.Feat{{ .Name }}Model.FindPageByWhereCount(l.ctx, whereStr)
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