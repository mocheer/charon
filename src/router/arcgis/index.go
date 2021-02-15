package arcgis

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/src/core/res"
	"github.com/mocheer/charon/src/models/types"
)

// Use 初始化 arcgis 路由
func Use(api fiber.Router) {
	router := api.Group("/arcgis")
	// GetCapabilities
	router.Get("/:name/capabilities", getCapabilities)
	// GetTile
	router.Get("/:name/tile/:z/:y/:x", getTile)
}

// getCapabilities
func getCapabilities(c *fiber.Ctx) error {
	nameParam := c.Params("name")
	server, err := NewTileServer(filepath.Join("./data/dmap", nameParam, "conf.xml"))
	if err == nil {
		return res.ResultOK(c, server)
	}
	return nil
}

// getTile
func getTile(c *fiber.Ctx) error {
	nameParam := c.Params("name")
	zParam := c.Params("z")
	yParam := c.Params("y")
	xParam := c.Params("x")
	server, err := NewTileServer(filepath.Join("./data/dmap", nameParam, "conf.xml"))
	if err == nil {
		x, _ := strconv.Atoi(xParam)
		y, _ := strconv.Atoi(yParam)
		z, _ := strconv.Atoi(zParam)
		data, err := server.ReadTile(types.Tile{
			Z: z, Y: y, X: x,
		})
		if err == nil {
			c.Type(strings.ToLower(server.FileFormat))
			return c.Send(data)
		}
	}

	return res.ResultError(c, 500, "读取瓦片错误", err)
}
