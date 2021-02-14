package gis

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/src/core/res"
	"github.com/mocheer/charon/src/models/types"
)

// Use 初始化 pipal 路由
func Use(api fiber.Router) {
	router := api.Group("/gis")
	// query
	router.Get("dmap/:name/:z/:y/:x", queryAppConfig)
}

// queryAppConfig
func queryAppConfig(c *fiber.Ctx) error {
	name := c.Params("name")
	z := c.Params("z")
	y := c.Params("y")
	x := c.Params("x")
	var Column, _ = strconv.Atoi(x)
	var Row, _ = strconv.Atoi(y)
	var Level, _ = strconv.Atoi(z)

	tc, err := NewMapServer(filepath.Join("D:/code-space/go/charon/data/dmap", name, "conf.xml"))
	if err == nil {
		data, tileErr := tc.ReadCompactTileV2(types.Tile{
			Row: Row, Level: Level, Column: Column,
		})

		if tileErr != nil {
			fmt.Println(tileErr)
		}
		c.Type(tc.FileFormat)

		return c.Send(data)
	}

	return res.ResultError(c, 500, "读取瓦片错误", err)

}
