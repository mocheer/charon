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
	router.Get("/:z/:y/:x", queryAppConfig)
}

// queryAppConfig
func queryAppConfig(c *fiber.Ctx) error {
	z := c.Params("z")
	y := c.Params("y")
	x := c.Params("x")
	var Column, _ = strconv.Atoi(x)
	var Row, _ = strconv.Atoi(y)
	var Level, _ = strconv.Atoi(z)
	fmt.Println(Level, Row, Column)
	// tc, err := NewEsri(filepath.Join("D:/code-space/go/charon/data/compactcacheV2", "conf.xml"))
	tc, err := NewEsri(filepath.Join("D:/code-space/go/charon/data/Layers", "conf.xml"))
	if err == nil {
		data, tileErr := tc.ReadCompactTileV2(types.Tile{
			Row: Row, Level: Level, Column: Column,
		})

		if tileErr != nil {
			fmt.Println(tileErr)
		}

		c.Type("png")
		c.Write(data)

		return nil
		// return c.SendStream(bytes.NewReader(data))
	}
	// x: 54624
	// y: 26923
	// z: 16
	return res.ResultError(c, 500, "读取瓦片错误", err)

}
