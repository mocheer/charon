package proxies

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/req"
)

// Use
func Use(api fiber.Router) {
	router := api.Group("proxy")
	router.Get("/", ProxyHandle)
	router.Post("/", ProxyHandle)
	router.Get("/amap/*", ProxyAMap)
}

func ProxyHandle(c *fiber.Ctx) error {
	args := &ProxyArgs{}
	req.MustParseArgs(c, args)
	return proxyURL(c, args.URL)
}
