package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DemoCurdModel = (*customDemoCurdModel)(nil)

type (
	// DemoCurdModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDemoCurdModel.
	DemoCurdModel interface {
		demoCurdModel
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
