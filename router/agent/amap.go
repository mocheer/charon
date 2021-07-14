package agent

import (
	"github.com/gofiber/fiber/v2"
)

// ProxyAMap 高德地图服务代理
func ProxyAMap(c *fiber.Ctx) error {
	url := "https://restapi.amap.com/" + c.Params("*") + "?key=77af3370f55d0399278ded758b023f59&" + c.Context().QueryArgs().String()
	return proxyURL(c, url)
}
