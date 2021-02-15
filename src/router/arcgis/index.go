package arcgis

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/src/core/res"
	"github.com/mocheer/charon/src/models/types"
)

// Use 初始化 pipal 路由
func Use(api fiber.Router) {
	router := api.Group("/arcgis")
	// query
	router.Get("/tile/:name/:z/:y/:x", getTile)
}

// getTile
func getTile(c *fiber.Ctx) error {
	name := c.Params("name")
	z := c.Params("z")
	y := c.Params("y")
	x := c.Params("x")
	var Column, _ = strconv.Atoi(x)
	var Row, _ = strconv.Atoi(y)
	var Level, _ = strconv.Atoi(z)

	tc, err := NewTileServer(filepath.Join("./data/dmap", name, "conf.xml"))
	if err == nil {
		data, err := tc.ReadTile(types.Tile{
			Row: Row, Level: Level, Column: Column,
		})
		if err == nil {
			c.Type(strings.ToLower(tc.FileFormat))
			return c.Send(data)
		}
	}

	return res.ResultError(c, 500, "读取瓦片错误", err)
}
