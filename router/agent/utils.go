package agent

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
	//
	// c.Request().Header.Reset() // 已知img标签不能直接reset
	// 重置host，改变源
	// @see http-proxy(nodejs) changeOrigin
	// Host 表示当前请求要被发送的目的地，仅包括域名和端口号
	// Origin CORS跨域请求中会有Origin，普通请求没有这个header
	c.Request().Header.SetHost(u.Host)
	c.Request().Header.SetRequestURI(u.RequestURI())
	// CSRF 防御之一：验证HTTP Referer 字段
	// Referer 记录了该HTTP请求的来源地址，在通常情况下，访问一个安全受限页面的请求必须来自于同一个网站，如果被检测出Referer是其他网站会判断为CSRF攻击，从而禁止访问
	// Referer 表示当前请求资源所在页面的完整路径：协议+域名+查询参数（不包含锚点信息） => 相当于页面的全路径
	c.Request().Header.SetReferer(u.String())
	//
	if err := proxy.Do(c, urlstr); err != nil {
		return err
	}
	// Remove Server header from response
	c.Response().Header.Del(fiber.HeaderServer)
	return nil
}
