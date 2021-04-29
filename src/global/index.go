package global

import (
	"github.com/mocheer/pluto/fs"
	"github.com/mocheer/pluto/logger"
)

// DB 数据库对象
var DB *DataBase

// Log 日志操作对象
var Log = logger.New()

// Config 全局配置
var Config *AppConfig

func Init() {
	// 初始化日志
	Log.Init()
	// 初始化配置
	Config = &AppConfig{
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
	// 读取应用配置
	err := fs.ReadJSON("config/app.json", Config)
	if err != nil {
		Log.Error("读取配置失败：", err)
	}
	// 连接数据库
	DB, err = openDB()
	if err != nil {
		Log.Error("数据库连接失败", err)
	}
	//
	Log.Info("启动成功")
}
