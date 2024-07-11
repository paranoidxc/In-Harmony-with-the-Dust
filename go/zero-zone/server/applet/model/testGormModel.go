package model

import (
	"gorm.io/gorm"
)

var _ TestGormModel = (*customTestGormModel)(nil)

type (
	// TestGormModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTestGormModel.
	TestGormModel interface {
		testGormModel
	}

	customTestGormModel struct {
		*defaultTestGormModel
	}
)

// NewTestGormModel returns a model for the database table.
func NewTestGormModel(gorm *gorm.DB) TestGormModel {
	return &customTestGormModel{
		defaultTestGormModel: newTestGormModel(gorm),
	}
}
