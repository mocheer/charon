package proxies

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

// proxyURL 常用于需要重写url+params的代理服务
func proxyURL(c *fiber.Ctx, urlstr string) error {
	u, err := url.Parse(urlstr)
	if err != nil {
		return err
	}
	// c.Request().Header.Reset();
	// 重置host，改变源
	// @see http-proxy(nodejs) changeOrigin
	c.Request().Header.SetHost(u.Host)
	//
	if err := proxy.Do(c, urlstr); err != nil {
		return err
	}
	// Remove Server header from response
	c.Response().Header.Del(fiber.HeaderServer)
	return nil
}
