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
	router.Get("/app/:name", queryStipule)
	router.Get("/page/:stipule/:name", store.GlobalCache, queryPetiole)
	router.Get("/view/:name", store.GlobalCache, queryBlade)
}

// queryStipule
func queryStipule(c *fiber.Ctx) error {
	name := c.Params("name")

	var stipule tables.Stipule
	global.Db.Where(&tables.Stipule{Name: name}).First(&stipule)
	return res.ResultOK(c, stipule)
}

// queryPetiole
func queryPetiole(c *fiber.Ctx) error {
	stipule := c.Params("stipule")
	name, _ := url.QueryUnescape(c.Params("name"))
	var petiole tables.Petiole
	global.Db.Where(&tables.Petiole{AppName: stipule, Name: name}).First(&petiole)
	return res.ResultOK(c, petiole)
}

// queryBlade
func queryBlade(c *fiber.Ctx) error {
	name := c.Params("name")
	var blade tables.Blade
	global.Db.Where(&tables.Blade{Name: name}).First(&blade)
	return res.ResultOK(c, blade)
}
