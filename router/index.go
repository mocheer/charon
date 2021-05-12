package router

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/mw"

	"github.com/mocheer/charon/router/arcgis"
	"github.com/mocheer/charon/router/auth"
	"github.com/mocheer/charon/router/dmap"
	"github.com/mocheer/charon/router/pipal"
	"github.com/mocheer/charon/router/proxies"
	"github.com/mocheer/charon/router/query"
	"github.com/mocheer/charon/router/structure"
	"github.com/mocheer/charon/router/upload"
)

// Init 初始化路由
func Init() {
	app := fiber.New()
	mw.Use(app)
	//
	api := app.Group("/api")
	v1 := api.Group("/v1")
	//
	auth.Use(v1)
	//
	mw.UseProtected(app)
	//
	pipal.Use(v1)
	query.Use(v1)
	upload.Use(v1)
	structure.Use(v1)
	proxies.Use(v1)
	arcgis.Use(v1)
	dmap.Use(v1)
	// 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("not found")
	})
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Fatal(app.Listen(global.Config.Port))
	// app.ListenTLS(":443", "./cert.pem", "./cert.key");//2.3.0
}
