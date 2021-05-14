package main

import (
	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/router"
)

var mode string

func main() {
	// 全局初始化
	global.Init(mode)
	// 初始化路由
	router.Init()
}
