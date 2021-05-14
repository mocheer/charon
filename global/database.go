package global

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 数据库对象
var DB *gorm.DB

// Open 连接数据库
func openDB() (db *gorm.DB, err error) {
	if Config.DSN == "" {
		return
	}

	db, err = gorm.Open(postgres.Open(Config.DSN), &gorm.Config{
		AllowGlobalUpdate: false,
		Logger:            logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(8)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(128)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	return
}
