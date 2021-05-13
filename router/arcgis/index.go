package arcgis

import (
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/model/types"
	"github.com/mocheer/charon/mw"
	"github.com/mocheer/charon/res"
)

// cacheServer
var cacheServer map[string]*TileServer = map[string]*TileServer{}

// Use 初始化 arcgis 路由
func Use(api fiber.Router) {
	router := api.Group("/arcgis")
	// GetCapabilities
	router.Get("/:name/capabilities", mw.Cache, getCapabilities)
	// GetTile
	router.Get("/:name/tile/:z/:y/:x", mw.CacheControl, getTile)
	//
}

// getCapabilities 获取元数据信息
func getCapabilities(c *fiber.Ctx) error {
	nameParam := c.Params("name")
	server, err := NewTileServer(filepath.Join(BaseDirectory, nameParam, "conf.xml"))
	if err == nil {
		return res.JSON(c, server)
	}
	return res.Error(c, "获取服务元数据错误", err)
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
	return res.Error(c, "读取瓦片错误", err)
}
