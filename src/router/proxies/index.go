package proxies

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

// Use
func Use(api fiber.Router) {
	router := api.Group("proxy")

	router.Use(proxy.Balancer(proxy.Config{
		Servers: []string{
			"http://t0.tianditu.gov.cn/",
		},
	}))
}
