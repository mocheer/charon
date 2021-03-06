package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/mw"
	"github.com/mocheer/charon/router/agent"
	"github.com/mocheer/charon/router/arcgis"
	"github.com/mocheer/charon/router/auth"
	"github.com/mocheer/charon/router/dmap"
	"github.com/mocheer/charon/router/pipal"
	"github.com/mocheer/charon/router/query"
	"github.com/mocheer/charon/router/reptile"
	"github.com/mocheer/charon/router/upload"
	"github.com/mocheer/charon/router/ws"
)

// v1_init /api/v1/xxx/xx
func v1_init(api fiber.Router) {
	v1 := api.Group("/v1")
	//
	auth.Use(v1)
	// 在auth之后执行token认证中间件
	mw.UseProtected(v1)
	//
	agent.Use(v1)
	pipal.Use(v1)
	query.Use(v1)
	upload.Use(v1)
	arcgis.Use(v1)
	dmap.Use(v1)
	ws.Use(v1)
	reptile.Use(v1)
}
