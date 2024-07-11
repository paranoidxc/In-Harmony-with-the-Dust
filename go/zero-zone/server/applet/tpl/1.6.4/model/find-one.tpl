func (m *default{{.upperStartCamelObject}}Model) FindOne(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error) {
	{{if .withCache}}{{.cacheKey}}
	var resp {{.upperStartCamelObject}}
	err := m.QueryRowCtx(ctx, &resp, {{.cacheKeyVariable}}, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query :=  fmt.Sprintf("select %s from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}} and `deleted_at` is null  limit 1", {{.lowerStartCamelObject}}Rows, m.table)
		return conn.QueryRowCtx(ctx, v, query, {{.lowerStartCamelPrimaryKey}})
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}{{else}}query := fmt.Sprintf("select %s from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}} limit 1", {{.lowerStartCamelObject}}Rows, m.table)
	var resp {{.upperStartCamelObject}}
	err := m.conn.QueryRowCtx(ctx, &resp, query, {{.lowerStartCamelPrimaryKey}})
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}{{end}}
}

func (m *default{{.upperStartCamelObject}}Model) FindAllByWhere(ctx context.Context, where string) ([]*{{.upperStartCamelObject}}, error) {
    query := fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY {{.originalPrimaryKey}} DESC", {{.lowerStartCamelObject}}Rows, m.table, where)
	var resp []*{{.upperStartCamelObject}}
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindAllByWhereCount(ctx context.Context, where string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT({{.originalPrimaryKey}}) FROM %s WHERE %s", m.table, where)
	var resp int64
	err := m.QueryRowNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindPageByWhere(ctx context.Context, where string, page int64, limit int64) ([]*{{.upperStartCamelObject}}, error) {
	offset := (page - 1) * limit
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY {{.originalPrimaryKey}} DESC LIMIT %d,%d", {{.lowerStartCamelObject}}Rows, m.table, where, offset, limit)
	var resp []*{{.upperStartCamelObject}}
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindPageByWhereCount(ctx context.Context, where string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT({{.originalPrimaryKey}}) FROM %s WHERE %s", m.table, where)
	var resp int64
	err := m.QueryRowNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}