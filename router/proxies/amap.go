package proxies

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

// ProxyAMap 高德地图服务代理
func ProxyAMap(c *fiber.Ctx) error {
	url := "https://restapi.amap.com/" + c.Params("*") + "?key=77af3370f55d0399278ded758b023f59&" + c.Context().QueryArgs().String()
	return proxyURL(c, url)
}

// proxyURL 常用于需要重写url+params的代理服务
func proxyURL(c *fiber.Ctx, url string) error {
	if err := proxy.Do(c, url); err != nil {
		return err
	}
	// Remove Server header from response
	c.Response().Header.Del(fiber.HeaderServer)
	return nil
}
