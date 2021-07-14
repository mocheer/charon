package agent

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/req"
)

// Use
func Use(api fiber.Router) {
	router := api.Group("agent")
	router.Get("/", ProxyHandle)
	router.Post("/", ProxyHandle)
	//
	router.Get("/raw/*", ProxyHandle2)
	router.Post("/raw/*", ProxyHandle2)
	// 高德地图
	router.Get("/amap/*", ProxyAMap)
	// 综合气象数据
	router.Get("/image.cma/*", ProxyImageCma)
}

// ProxyHandle
func ProxyHandle(c *fiber.Ctx) error {
	args := &ProxyArgs{}
	req.MustParseArgs(c, args)
	//
	return proxyURL(c, args.URL)
}

func ProxyHandle2(c *fiber.Ctx) error {
	return proxyURL(c, "http://"+c.Params("*")+"?"+c.Context().QueryArgs().String())
}
