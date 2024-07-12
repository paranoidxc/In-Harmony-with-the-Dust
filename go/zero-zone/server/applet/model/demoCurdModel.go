package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DemoCurdModel = (*customDemoCurdModel)(nil)

type (
	// DemoCurdModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDemoCurdModel.
	DemoCurdModel interface {
		demoCurdModel

		FindAllByWhere(ctx context.Context, where string) ([]*DemoCurd, error)
		FindAllByWhereCount(ctx context.Context, where string) (int64, error)
		FindPageByWhere(ctx context.Context, where string, page int64, limit int64) ([]*DemoCurd, error)
		FindPageByWhereCount(ctx context.Context, where string) (int64, error)
	}

	customDemoCurdModel struct {
		*defaultDemoCurdModel
	}
)

// NewDemoCurdModel returns a model for the database table.
func NewDemoCurdModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) DemoCurdModel {
	return &customDemoCurdModel{
		defaultDemoCurdModel: newDemoCurdModel(conn, c, opts...),
	}
}

func (m *defaultDemoCurdModel) FindAllByWhere(ctx context.Context, where string) ([]*DemoCurd, error) {
	query := fmt.Sprintf("SELECT %s FROM %s AS t WHERE %s ORDER BY t.id DESC", demoCurdRows, m.table, where)
	var resp []*DemoCurd
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultDemoCurdModel) FindAllByWhereCount(ctx context.Context, where string) (int64, error) {
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

func (m *defaultDemoCurdModel) FindPageByWhere(ctx context.Context, where string, page int64, limit int64) ([]*DemoCurd, error) {
	offset := (page - 1) * limit
	query := fmt.Sprintf("SELECT %s FROM %s AS t WHERE %s ORDER BY t.id DESC LIMIT %d,%d", demoCurdRows, m.table, where, offset, limit)
	var resp []*DemoCurd
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultDemoCurdModel) FindPageByWhereCount(ctx context.Context, where string) (int64, error) {
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
