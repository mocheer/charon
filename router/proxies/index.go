package proxies

import (
	"github.com/gofiber/fiber/v2"
)

// Use
func Use(api fiber.Router) {
	router := api.Group("proxy")
	router.Get("/amap/*", ProxyAMap)
}
