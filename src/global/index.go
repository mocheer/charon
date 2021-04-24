package global

import (
	"gorm.io/gorm"
)

// Db 数据库对象
var Db *gorm.DB

// Config 全局配置
var Config = &AppConfig{
	Name: "charon",
	Mode: "production",
	Port: ":9212",
	Static: map[string]StaticConfig{
		"/": {
			Mode: "history",
			Dir:  "./public",
		},
	},
}
