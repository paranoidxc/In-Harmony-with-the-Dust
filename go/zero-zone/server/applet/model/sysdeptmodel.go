package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysDeptModel = (*customSysDeptModel)(nil)

type (
	// SysDeptModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysDeptModel.
	SysDeptModel interface {
		sysDeptModel
		FindAll(ctx context.Context) ([]*SysDept, error)
		FindSubDept(ctx context.Context, id int64) ([]*SysDept, error)
		FindEnable(ctx context.Context) ([]*SysDept, error)
		FindPagination(ctx context.Context, where string, page int64, limit int64) ([]*SysDept, error)
		FindPage(ctx context.Context, page int64, limit int64) ([]*SysDept, error)
		FindCount(ctx context.Context) (int64, error)
		FindPaginationCount(ctx context.Context, where string) (int64, error)
	}

	customSysDeptModel struct {
		*defaultSysDeptModel
	}
)

// NewSysDeptModel returns a model for the database table.
func NewSysDeptModel(conn sqlx.SqlConn, c cache.CacheConf) SysDeptModel {
	return &customSysDeptModel{
		defaultSysDeptModel: newSysDeptModel(conn, c),
	}
}

func (m *customSysDeptModel) FindPagination(ctx context.Context, where string, page int64, limit int64) ([]*SysDept, error) {
	query := ""
	offset := (page - 1) * limit
	if len(where) == 0 {
		query = fmt.Sprintf("SELECT %s FROM %s ORDER BY order_num DESC LIMIT %d, %d ", sysDeptRows, m.table, offset, limit)
	} else {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY order_num DESC LIMIT %d, %d ", sysDeptRows, m.table, where, offset, limit)
	}

	var resp []*SysDept
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customSysDeptModel) FindAll(ctx context.Context) ([]*SysDept, error) {
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY order_num DESC", sysDeptRows, m.table)
	var resp []*SysDept
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customSysDeptModel) FindEnable(ctx context.Context) ([]*SysDept, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE status=1 ORDER BY order_num DESC", sysDeptRows, m.table)
	var resp []*SysDept
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customSysDeptModel) FindSubDept(ctx context.Context, id int64) ([]*SysDept, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `parent_id`=%d", sysDeptRows, m.table, id)
	var resp []*SysDept
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customSysDeptModel) FindPage(ctx context.Context, page int64, limit int64) ([]*SysDept, error) {
	offset := (page - 1) * limit
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY order_num DESC LIMIT %d,%d", sysDeptRows, m.table, offset, limit)
	var resp []*SysDept
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customSysDeptModel) FindCount(ctx context.Context) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(id) FROM %s", m.table)
	var resp int64
	err := m.QueryRowNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *customSysDeptModel) FindPaginationCount(ctx context.Context, where string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE %s ", m.table, where)
	var resp int64
	err := m.QueryRowNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}
