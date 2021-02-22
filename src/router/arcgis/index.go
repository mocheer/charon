package arcgis

import (
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/src/core/res"
	"github.com/mocheer/charon/src/models/types"
	"github.com/mocheer/charon/src/router/store"
)

// cacheServer
var cacheServer map[string]*TileServer = map[string]*TileServer{}

// Use 初始化 arcgis 路由
func Use(api fiber.Router) {
	router := api.Group("/arcgis")
	// GetCapabilities
	router.Get("/:name/capabilities", store.GlobalCache, getCapabilities)
	// GetTile
	router.Get("/:name/tile/:z/:y/:x", store.NewCache(time.Hour*24*30), getTile)
	//
}

// getCapabilities 获取元数据信息
func getCapabilities(c *fiber.Ctx) error {
	nameParam := c.Params("name")
	server, err := NewTileServer(filepath.Join(BaseDirectory, nameParam, "conf.xml"))
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
	server, hasCache := cacheServer[nameParam]
	var err error
	if !hasCache {
		server, err = NewTileServer(filepath.Join(BaseDirectory, nameParam, "conf.xml"))
		cacheServer[nameParam] = server
	}
	if err == nil {
		x, _ := strconv.Atoi(xParam)
		y, _ := strconv.Atoi(yParam)
		z, _ := strconv.Atoi(zParam)
		data, err := server.ReadTile(types.Tile{
			Z: z, Y: y, X: x,
		})
		if err == nil {
			c.Type(server.TileFormat)
			return c.Send(data)
		}
	}
	//
	// return res.ResultError(c, 500, "读取瓦片错误", err)
	return nil
}
