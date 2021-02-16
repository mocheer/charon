package dmap

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/src/core/res"
)

// Use 初始化 arcgis 路由
func Use(api fiber.Router) {
	router := api.Group("/dmap")
	//
	router.Get("/dynamic-layer", getDynamicLayer)
}

//
// getDynamicLayer
func getDynamicLayer(c *fiber.Ctx) error {
	//
	return res.ResultError(c, 500, "获取数据错误", nil)
}
