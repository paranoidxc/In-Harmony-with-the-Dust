package {{.pkg}}
{{if .withCache}}
import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)
{{else}}

import "github.com/zeromicro/go-zero/core/stores/sqlx"
{{end}}
var _ {{.upperStartCamelObject}}Model = (*custom{{.upperStartCamelObject}}Model)(nil)

type (
	// {{.upperStartCamelObject}}Model is an interface to be customized, add more methods here,
	// and implement the added methods in custom{{.upperStartCamelObject}}Model.
	{{.upperStartCamelObject}}Model interface {
		{{.lowerStartCamelObject}}Model
		{{if not .withCache}}withSession(session sqlx.Session) {{.upperStartCamelObject}}Model{{end}}
		FindAllByWhere(ctx context.Context, where string) ([]*{{.upperStartCamelObject}}, error)
        FindAllByWhereCount(ctx context.Context, where string) (int64, error)
        FindPageByWhere(ctx context.Context, where string, page int64, limit int64) ([]*{{.upperStartCamelObject}}, error)
        FindPageByWhereCount(ctx context.Context, where string) (int64, error)
	}

	custom{{.upperStartCamelObject}}Model struct {
		*default{{.upperStartCamelObject}}Model
	}
)

// New{{.upperStartCamelObject}}Model returns a model for the database table.
func New{{.upperStartCamelObject}}Model(conn sqlx.SqlConn{{if .withCache}}, c cache.CacheConf, opts ...cache.Option{{end}}) {{.upperStartCamelObject}}Model {
	return &custom{{.upperStartCamelObject}}Model{
		default{{.upperStartCamelObject}}Model: new{{.upperStartCamelObject}}Model(conn{{if .withCache}}, c, opts...{{end}}),
	}
}

{{if not .withCache}}
func (m *custom{{.upperStartCamelObject}}Model) withSession(session sqlx.Session) {{.upperStartCamelObject}}Model {
    return New{{.upperStartCamelObject}}Model(sqlx.NewSqlConnFromSession(session))
}
{{end}}


func (m *default{{.upperStartCamelObject}}Model) FindAllByWhere(ctx context.Context, where string) ([]*{{.upperStartCamelObject}}, error) {
    query := fmt.Sprintf("SELECT %s FROM %s AS t WHERE %s ORDER BY t.id DESC", {{.lowerStartCamelObject}}Rows, m.table, where)
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
	query := fmt.Sprintf("SELECT COUNT(t.id) FROM %s AS t WHERE %s", m.table, where)
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
	query := fmt.Sprintf("SELECT %s FROM %s AS t WHERE %s ORDER BY t.id DESC LIMIT %d,%d", {{.lowerStartCamelObject}}Rows, m.table, where, offset, limit)
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
	query := fmt.Sprintf("SELECT COUNT(t.id) FROM %s AS t WHERE %s", m.table, where)
	var resp int64
	err := m.QueryRowNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}
