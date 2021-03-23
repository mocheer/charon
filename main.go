package main

import (
	_ "net/http/pprof"

	"github.com/mocheer/charon/src/core/db"
	"github.com/mocheer/charon/src/global"
	"github.com/mocheer/charon/src/logger"
	"github.com/mocheer/charon/src/router"
)

func main() {
	logger.Init()
	// 读取应用配置
	err := global.Config.ReadJSON("./config/app.json")
	if err != nil {
		logger.Error("读取配置文件失败：", err)
	}
	// 连接数据库
	if global.Config.DbDSN != "" {
		global.Db = db.Open(global.Config.DbDSN)
	}
	// 初始化服务
	router.Init()
	//
	logger.Info("启动成功")

}
