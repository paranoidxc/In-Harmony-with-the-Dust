package database

import (
	"time"
	"zero-zone/applet/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DBConn struct {
	ConnGorm *gorm.DB
}

// NewDB  连接并初始化数据库
func NewDB(dataSource string) *DBConn {
	db, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   false,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: logger.Default.LogMode(logger.Info),
		// 启用更新时间戳功能
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		panic("连接数据库失败")
	}

	d := &DBConn{
		ConnGorm: db,
	}
	InitDB(db)
	return d
}
func InitDB(db *gorm.DB) {
	if err := db.AutoMigrate(
		&model.TestGorm{},
	); err != nil {
		panic(err)
	}
}
