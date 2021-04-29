package main

import (
	_ "net/http/pprof"

	"github.com/mocheer/charon/src/global"
	"github.com/mocheer/charon/src/router"
)

func main() {
	// 全局初始化
	global.Init()
	// 初始化路由
	router.Init()
}
