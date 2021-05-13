package router

import (
	"github.com/gofiber/fiber/v2"
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

// apiV1 /api/v1/xxx/xx
func apiV1(api fiber.Router) {
	v1 := api.Group("/v1")
	//
	auth.Use(v1)
	// 在auth之后执行token认证中间件
	mw.UseProtected(v1)
	//
	pipal.Use(v1)
	query.Use(v1)
	upload.Use(v1)
	structure.Use(v1)
	proxies.Use(v1)
	arcgis.Use(v1)
	dmap.Use(v1)
}
