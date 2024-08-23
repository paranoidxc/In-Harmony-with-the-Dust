package svc

import (
	"zero-zone/applet/api/internal/config"
	"zero-zone/applet/api/internal/middleware"
	"zero-zone/applet/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config              config.Config
	Redis               *redis.Redis
	PermMenuAuth        rest.Middleware
	ReqRespLog          rest.Middleware
	LogActionMiddleware rest.Middleware
	SysUserModel        model.SysUserModel
	SysPermMenuModel    model.SysPermMenuModel
	SysRoleModel        model.SysRoleModel
	SysDeptModel        model.SysDeptModel
	SysJobModel         model.SysJobModel
	SysProfessionModel  model.SysProfessionModel
	SysDictionaryModel  model.SysDictionaryModel
	SysLogModel         model.SysLogModel
	FeatSysRegionModel  model.SysRegionModel
	FeatDemoCurdModel   model.DemoCurdModel
	FeatTestGormModel   model.TestGormModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	//mysqlGorm := database.NewDB(c.Mysql.DataSource) //引入数据库得到数据库链接
	redisClient := redis.New(c.Redis.Host, func(r *redis.Redis) {
		r.Type = c.Redis.Type
		r.Pass = c.Redis.Pass
	})
	//原始方式
	// return &ServiceContext{
	// 	Config:             c,
	// 	Redis:              redisClient,
	// 	PermMenuAuth:       middleware.NewPermMenuAuthMiddleware(redisClient).Handle,
	// 	SysUserModel:       model.NewSysUserModel(mysqlConn, c.Cache),
	// 	SysPermMenuModel:   model.NewSysPermMenuModel(mysqlConn, c.Cache),
	// 	SysRoleModel:       model.NewSysRoleModel(mysqlConn, c.Cache),
	// 	SysDeptModel:       model.NewSysDeptModel(mysqlConn, c.Cache),
	// 	SysJobModel:        model.NewSysJobModel(mysqlConn, c.Cache),
	// 	SysProfessionModel: model.NewSysProfessionModel(mysqlConn, c.Cache),
	// 	SysDictionaryModel: model.NewSysDictionaryModel(mysqlConn, c.Cache),
	// 	SysLogModel:        model.NewSysLogModel(mysqlConn, c.Cache),
	// 	FeatTdFirmModel:    model.NewTdFirmModel(mysqlConn, c.Cache),
	// }
	//gorm方式
	return &ServiceContext{
		Config:              c,
		Redis:               redisClient,
		PermMenuAuth:        middleware.NewPermMenuAuthMiddleware(redisClient).Handle,
		LogActionMiddleware: middleware.NewLogActionMiddleware().Handle,
		SysUserModel:        model.NewSysUserModel(mysqlConn, c.Cache),
		SysPermMenuModel:    model.NewSysPermMenuModel(mysqlConn, c.Cache),
		SysRoleModel:        model.NewSysRoleModel(mysqlConn, c.Cache),
		SysDeptModel:        model.NewSysDeptModel(mysqlConn, c.Cache),
		SysJobModel:         model.NewSysJobModel(mysqlConn, c.Cache),
		SysProfessionModel:  model.NewSysProfessionModel(mysqlConn, c.Cache),
		SysDictionaryModel:  model.NewSysDictionaryModel(mysqlConn, c.Cache),
		SysLogModel:         model.NewSysLogModel(mysqlConn, c.Cache),
		FeatSysRegionModel:  model.NewSysRegionModel(mysqlConn, c.Cache),
		FeatDemoCurdModel:   model.NewDemoCurdModel(mysqlConn, c.Cache),
		//FeatTestGormModel:   model.NewTestGormModel(mysqlGorm.ConnGorm),
	}
}
