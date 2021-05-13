package router

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/mw"
)

// Init 初始化路由
func Init() {
	app := fiber.New()
	mw.Use(app)
	//
	api := app.Group("/api")
	apiV1(api)
	// 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("not found")
	})
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Fatal(app.Listen(global.Config.Port))
	// app.ListenTLS(":443", "./cert.pem", "./cert.key");//2.3.0
}
