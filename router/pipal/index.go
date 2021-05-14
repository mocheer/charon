package pipal

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/model/tables"
	"github.com/mocheer/charon/mw"
	"github.com/mocheer/charon/res"
)

// Use 初始化 pipal 路由
func Use(api fiber.Router) {
	router := api.Group("/pipal")
	// query
	router.Get("/app/:name", queryAppConfig)
	router.Get("/page/:appName/:name", mw.Cache, queryPageConfig)
	router.Get("/view/:name", mw.Cache, queryViewConfig)
}

// queryAppConfig
func queryAppConfig(c *fiber.Ctx) error {
	name := c.Params("name")
	//
	var appConfig tables.AppConfig
	global.DB.Where(&tables.AppConfig{Name: name}).FirstOrCreate(&appConfig, &tables.AppConfig{Name: name, Enabled: true})
	return res.JSON(c, appConfig)
}

// queryPageConfig
func queryPageConfig(c *fiber.Ctx) error {
	var appName = c.Params("appName")
	var name, _ = url.QueryUnescape(c.Params("name"))
	var petiole tables.PageConfig
	global.DB.Where(&tables.PageConfig{AppName: appName, Name: name}).FirstOrCreate(&petiole, &tables.PageConfig{AppName: appName, Name: name})
	return res.JSON(c, petiole)
}

// queryViewConfig
func queryViewConfig(c *fiber.Ctx) error {
	name := c.Params("name")
	var viewConfig tables.ViewConfig
	global.DB.Where(&tables.ViewConfig{Name: name}).First(&viewConfig)
	return res.JSON(c, viewConfig)
}
