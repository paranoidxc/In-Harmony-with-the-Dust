// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	sysRegionFieldNames          = builder.RawFieldNames(&SysRegion{})
	sysRegionRows                = strings.Join(sysRegionFieldNames, ",")
	sysRegionRowsExpectAutoSet   = strings.Join(stringx.Remove(sysRegionFieldNames, "`id`", "`created_at`", "`deleted_at`", "`updated_at`"), ",")
	sysRegionRowsWithPlaceHolder = strings.Join(stringx.Remove(sysRegionFieldNames, "`id`", "`created_at`", "`deleted_at`", "`updated_at`"), "=?,") + "=?"

	cacheZeroZoneSysRegionIdPrefix = "cache:zeroZone:sysRegion:id:"
	cacheZeroZoneSysRegionNoPrefix = "cache:zeroZone:sysRegion:no:"
)

type (
	sysRegionModel interface {
		Insert(ctx context.Context, data *SysRegion) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*SysRegion, error)
		FindOneByNo(ctx context.Context, no int64) (*SysRegion, error)
		Update(ctx context.Context, data *SysRegion) error
		Delete(ctx context.Context, id int64) error
		Deletes(ctx context.Context, ids []int64) error
	}

	defaultSysRegionModel struct {
		sqlc.CachedConn
		table string
	}

	SysRegion struct {
		Id        int64     `db:"id"`        // ID
		No        int64     `db:"no"`        // 区域编码
		Name      string    `db:"name"`      // 区域名称
		ParentNo  int64     `db:"parent_no"` // 上级区域
		Code      string    `db:"code"`      // 电话区号
		Level     int64     `db:"level"`     // 区域级别
		Typename  string    `db:"typename"`  // 级别名称
		PySzm     string    `db:"py_szm"`    // 首字母
		IsDel     int64     `db:"is_del"`
		CreatedAt time.Time `db:"created_at"` // 创建时间
		UpdatedAt time.Time `db:"updated_at"` // 更新时间
		DeletedAt string    `db:"deleted_at"` // 删除时间
	}
)

func newSysRegionModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultSysRegionModel {
	return &defaultSysRegionModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`sys_region`",
	}
}

func (m *defaultSysRegionModel) Delete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	zeroZoneSysRegionIdKey := fmt.Sprintf("%s%v", cacheZeroZoneSysRegionIdPrefix, id)
	zeroZoneSysRegionNoKey := fmt.Sprintf("%s%v", cacheZeroZoneSysRegionNoPrefix, data.No)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		//query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		query := fmt.Sprintf("update %s set `deleted_at` = now(), is_del = 1 where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, zeroZoneSysRegionIdKey, zeroZoneSysRegionNoKey)
	return err
}

func (m *defaultSysRegionModel) Deletes(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		data, err := m.FindOne(ctx, id)
		if err != nil {
			return err
		}

		zeroZoneSysRegionIdKey := fmt.Sprintf("%s%v", cacheZeroZoneSysRegionIdPrefix, id)
		zeroZoneSysRegionNoKey := fmt.Sprintf("%s%v", cacheZeroZoneSysRegionNoPrefix, data.No)
		_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
			//query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
			query := fmt.Sprintf("update %s set `deleted_at` = now(), is_del = 1 where `id` = ?", m.table)
			return conn.ExecCtx(ctx, query, id)
		}, zeroZoneSysRegionIdKey, zeroZoneSysRegionNoKey)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *defaultSysRegionModel) FindOne(ctx context.Context, id int64) (*SysRegion, error) {
	zeroZoneSysRegionIdKey := fmt.Sprintf("%s%v", cacheZeroZoneSysRegionIdPrefix, id)
	var resp SysRegion
	err := m.QueryRowCtx(ctx, &resp, zeroZoneSysRegionIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? and is_del = 0  limit 1", sysRegionRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultSysRegionModel) FindOneByNo(ctx context.Context, no int64) (*SysRegion, error) {
	zeroZoneSysRegionNoKey := fmt.Sprintf("%s%v", cacheZeroZoneSysRegionNoPrefix, no)
	var resp SysRegion
	err := m.QueryRowIndexCtx(ctx, &resp, zeroZoneSysRegionNoKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where `no` = ? limit 1", sysRegionRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, no); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultSysRegionModel) Insert(ctx context.Context, data *SysRegion) (sql.Result, error) {
	zeroZoneSysRegionIdKey := fmt.Sprintf("%s%v", cacheZeroZoneSysRegionIdPrefix, data.Id)
	zeroZoneSysRegionNoKey := fmt.Sprintf("%s%v", cacheZeroZoneSysRegionNoPrefix, data.No)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, sysRegionRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.No, data.Name, data.ParentNo, data.Code, data.Level, data.Typename, data.PySzm, data.IsDel)
	}, zeroZoneSysRegionIdKey, zeroZoneSysRegionNoKey)
	return ret, err
}

func (m *defaultSysRegionModel) Update(ctx context.Context, newData *SysRegion) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	zeroZoneSysRegionIdKey := fmt.Sprintf("%s%v", cacheZeroZoneSysRegionIdPrefix, data.Id)
	zeroZoneSysRegionNoKey := fmt.Sprintf("%s%v", cacheZeroZoneSysRegionNoPrefix, data.No)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, sysRegionRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.No, newData.Name, newData.ParentNo, newData.Code, newData.Level, newData.Typename, newData.PySzm, newData.IsDel, newData.Id)
	}, zeroZoneSysRegionIdKey, zeroZoneSysRegionNoKey)
	return err
}

func (m *defaultSysRegionModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheZeroZoneSysRegionIdPrefix, primary)
}

func (m *defaultSysRegionModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? and `deleted_at` is null limit 1", sysRegionRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultSysRegionModel) tableName() string {
	return m.table
}
