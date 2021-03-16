package pipal

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/src/core/res"
	"github.com/mocheer/charon/src/global"
	"github.com/mocheer/charon/src/models/tables"
	"github.com/mocheer/charon/src/router/store"
)

// Use 初始化 pipal 路由
func Use(api fiber.Router) {
	router := api.Group("/pipal")
	// query
	router.Get("/app/:name", queryAppConfig)
	router.Get("/page/:appName/:name", store.GlobalCache, queryPageConfig)
	router.Get("/view/:name", store.GlobalCache, queryViewConfig)
}

// queryAppConfig
func queryAppConfig(c *fiber.Ctx) error {
	name := c.Params("name")
	//
	var appConfig tables.AppConfig
	global.Db.Where(&tables.AppConfig{Name: name}).FirstOrCreate(&appConfig, &tables.AppConfig{Name: name, Enabled: true})
	return res.ResultOK(c, appConfig)
}

// queryPageConfig
func queryPageConfig(c *fiber.Ctx) error {
	var appName = c.Params("appName")
	var name, _ = url.QueryUnescape(c.Params("name"))
	var petiole tables.PageConfig
	global.Db.Where(&tables.PageConfig{AppName: appName, Name: name}).FirstOrCreate(&petiole, &tables.PageConfig{AppName: appName, Name: name})
	return res.ResultOK(c, petiole)
}

// queryViewConfig
func queryViewConfig(c *fiber.Ctx) error {
	name := c.Params("name")
	var viewConfig tables.ViewConfig
	global.Db.Where(&tables.ViewConfig{Name: name}).First(&viewConfig)
	return res.ResultOK(c, viewConfig)
}
