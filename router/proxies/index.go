package proxies

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

// Use 77af3370f55d0399278ded758b023f59
func Use(api fiber.Router) {
	router := api.Group("proxy")
	tdt := router.Group("amap")
	tdt.Use(proxy.Balancer(proxy.Config{
		Servers: []string{
			"",
		},
	}))
}
