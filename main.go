package main

import (
	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/router"
)

func main() {
	// 全局初始化
	global.Init()
	// 初始化路由
	router.Init()
}
