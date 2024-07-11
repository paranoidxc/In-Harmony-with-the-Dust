
func (l *{{ .Name }}ListLogic) {{ .Name }}List(req *types.{{ .Name }}ListReq) (resp *types.{{ .Name }}ListResp, err error) {
    where := " deleted_at IS NULL "
    /*
    {{- range $i, $v := .VueFields }}
    if len(strings.TrimSpace(req.{{ $v.Name }})) > 0 {
        where = where + fmt.Sprintf(" AND {{ $v.Column }} LIKE '%s'", "%"+strings.TrimSpace(req.{{ $v.Name }})+"%")
    }
    {{- end }}
    */
    feat{{ .Name }}List, err := l.svcCtx.Feat{{ .Name }}Model.FindAllByWhere(l.ctx, where)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var item types.{{ .Name }}
	{{ .Name }}List := make([]types.{{ .Name }}, 0)
	for _, v := range feat{{ .Name }}List {
		err := copier.Copy(&item, &v)
		item.CreatedAt = utils.Time2Str(v.CreatedAt)
       	item.UpdatedAt = utils.Time2Str(v.UpdatedAt)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		{{ .Name }}List = append({{ .Name }}List, item)
	}

	return &types.{{ .Name }}ListResp{
		List: {{ .Name }}List,
	}, nil
}