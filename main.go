package main

import (
	"fmt"

	"github.com/mocheer/charon/src/core/db"
	"github.com/mocheer/charon/src/global"
	"github.com/mocheer/charon/src/router"
)

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9912
// @BasePath /
func main() {
	// 读取应用配置
	err := global.Config.ReadJSON("./config/app.json")
	if err != nil {
		fmt.Println("读取配置文件失败：", err)
	}
	// 连接数据库
	if global.Config.DbDSN != "" {
		global.Db = db.Open(global.Config.DbDSN)
	}
	// 初始化服务
	router.Init()
}
