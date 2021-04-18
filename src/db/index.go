package db

import (
	"time"

	"github.com/mocheer/charon/src/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Open 连接数据库
func Open(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{AllowGlobalUpdate: false})
	if err != nil {
		logger.Error("数据库连接失败")
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("数据库连接失败", err)
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(8)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(128)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	//
	return db
}
