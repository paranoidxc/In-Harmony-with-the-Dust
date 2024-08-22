package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysRegionModel = (*customSysRegionModel)(nil)

type (
	// SysRegionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysRegionModel.
	SysRegionModel interface {
		sysRegionModel

		FindAllByWhere(ctx context.Context, where string) ([]*SysRegion, error)
		FindAllByWhereCount(ctx context.Context, where string) (int64, error)
		FindPageByWhere(ctx context.Context, where string, page int64, limit int64) ([]*SysRegion, error)
		FindPageByWhereCount(ctx context.Context, where string) (int64, error)
	}

	customSysRegionModel struct {
		*defaultSysRegionModel
	}
)

// NewSysRegionModel returns a model for the database table.
func NewSysRegionModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SysRegionModel {
	return &customSysRegionModel{
		defaultSysRegionModel: newSysRegionModel(conn, c, opts...),
	}
}

func (m *defaultSysRegionModel) FindAllByWhere(ctx context.Context, where string) ([]*SysRegion, error) {
	query := fmt.Sprintf("SELECT %s FROM %s AS t WHERE %s ORDER BY t.id DESC", sysRegionRows, m.table, where)
	var resp []*SysRegion
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultSysRegionModel) FindAllByWhereCount(ctx context.Context, where string) (int64, error) {
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

func (m *defaultSysRegionModel) FindPageByWhere(ctx context.Context, where string, page int64, limit int64) ([]*SysRegion, error) {
	offset := (page - 1) * limit
	query := fmt.Sprintf("SELECT %s FROM %s AS t WHERE %s ORDER BY t.id DESC LIMIT %d,%d", sysRegionRows, m.table, where, offset, limit)
	var resp []*SysRegion
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultSysRegionModel) FindPageByWhereCount(ctx context.Context, where string) (int64, error) {
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
