package router

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/mocheer/charon/cts"
	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/mw"
)

// Init 初始化路由
func Init() {
	app := fiber.New()
	mw.Use(app)
	//
	api := app.Group("/api")
	v1_init(api)
	// 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("not found")
	})
	// 监听
	var err error
	if global.Config.IsHTTPS {
		err = app.ListenTLS(global.Config.Port, cts.Cert_PEM, cts.Cert_KEY)
	} else {
		err = app.Listen(global.Config.Port)
	}
	log.Fatal(err)
}
