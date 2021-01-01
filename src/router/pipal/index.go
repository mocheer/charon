package pipal

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/src/core/res"
	"github.com/mocheer/charon/src/global"
	"github.com/mocheer/charon/src/models/tables"
	"github.com/mocheer/charon/src/router/auth"
	"github.com/mocheer/charon/src/router/store"
)

// Use 初始化 pipal 路由
func Use(api fiber.Router) {
	router := api.Group("/pipal")
	// query
	router.Get("/stipule/:name", queryStipule)
	router.Get("/petiole/:stipule/:name", store.GlobalCache, queryPetiole)
	router.Get("/blade/:name", store.GlobalCache, queryBlade)
	// update
	router.Post("/stipule", auth.GlobalProtected, updateStipule)
	router.Post("/petiole", auth.GlobalProtected, updatePetiole)
	router.Post("/blade", auth.GlobalProtected, updateBlade)
	// delete
	router.Delete("/stipule/:name", auth.GlobalProtected, deletePetiole)
	router.Delete("/petiole/:name", auth.GlobalProtected, deletePetiole)
	router.Delete("/blade/:name", auth.GlobalProtected, deleteBlade)
}

// queryStipule
func queryStipule(c *fiber.Ctx) error {
	name := c.Params("name")

	var stipule tables.Stipule
	global.Db.Where(&tables.Stipule{Name: name}).First(&stipule)
	return res.ResultOK(c, stipule)
}

// updateStipule
func updateStipule(c *fiber.Ctx) error {
	var stipule tables.Stipule
	if err := c.BodyParser(stipule); err != nil {

	}
	global.Db.Save(&stipule)
	return res.ResultOK(c, stipule)
}

// queryPetiole
func queryPetiole(c *fiber.Ctx) error {
	stipule := c.Params("stipule")
	name, _ := url.QueryUnescape(c.Params("name"))
	var petiole tables.Petiole
	global.Db.Where(&tables.Petiole{Stipule: stipule, Name: name}).First(&petiole)
	return res.ResultOK(c, petiole)
}

// updatePetiole
func updatePetiole(c *fiber.Ctx) error {
	var petiole tables.Petiole
	if err := c.BodyParser(petiole); err != nil {

	}
	global.Db.Save(&petiole)
	return res.ResultOK(c, petiole)
}

// queryBlade
func queryBlade(c *fiber.Ctx) error {
	name := c.Params("name")
	var blade tables.Blade
	global.Db.Where(&tables.Blade{Name: name}).First(&blade)
	return res.ResultOK(c, blade)
}

// updateBlade
func updateBlade(c *fiber.Ctx) error {
	var blade tables.Blade
	if err := c.BodyParser(blade); err != nil {

	}
	global.Db.Save(&blade)
	return res.ResultOK(c, blade)
}

// deleteBlade
func deleteBlade(c *fiber.Ctx) error {
	var blade tables.Blade
	if err := c.BodyParser(blade); err != nil {

	}
	global.Db.Delete(blade)
	return res.ResultOK(c, true)
}

// deletePetiole
func deletePetiole(c *fiber.Ctx) error {
	var blade tables.Blade
	if err := c.BodyParser(blade); err != nil {

	}
	global.Db.Delete(blade)
	return res.ResultOK(c, true)
}
