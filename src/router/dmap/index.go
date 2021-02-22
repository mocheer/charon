package dmap

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/src/core/res"
	"github.com/mocheer/charon/src/models/types"
)

// Use 初始化 dmap 路由
func Use(api fiber.Router) {
	router := api.Group("/dmap")
	//
	router.Get("/:id/image/:z/:y:x", getDynamicLayer)
}

// getDynamicLayer
func getDynamicLayer(c *fiber.Ctx) error {
	//
	idParam := c.Params("id")
	zParam := c.Params("z")
	yParam := c.Params("y")
	xParam := c.Params("x")

	id, _ := strconv.Atoi(idParam)
	x, _ := strconv.Atoi(xParam)
	y, _ := strconv.Atoi(yParam)
	z, _ := strconv.Atoi(zParam)

	layer := NewDynamicLayer(id, &types.Tile{
		Z: z, Y: y, X: x,
	})

	if layer != nil {

	}

	return res.ResultError(c, 500, "获取数据错误", nil)
}
