package mw

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/mocheer/pluto/ts/clock"
)

// Cache 全局缓存
var Cache = NewCache(time.Hour * 24)

// CacheControl 全局的前端强缓存
var CacheControl = NewCacheControl(clock.Month)

// NewCache 创建缓存
func NewCache(exp time.Duration) func(*fiber.Ctx) error {
	return cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration:   exp,
		CacheControl: true,
		KeyGenerator: func(c *fiber.Ctx) string {
			// c.OriginalURL() 在本地运行时有问题，因为用的是 UnsafeString 即 *(*string)(unsafe.Pointer(&bytes))) 来转换[]byte，高性能但不安全（用的同一个内存块，即使内容改变了，但没发现）
			// 但为什么远程部署后没问题?
			// 注意，发生的场景，是相同path，相同长度的url，只是参数不同（参数个数要相同）
			// /api/v1/query/s/petiole/first?r=11212654&where=%7B%22stipule%22:%22web%22,%22name%22:%22home%22%7D
			// /api/v1/query/s/petiole/first?r=11212654&where=%7B%22stipule%22:%22web%22,%22name%22:%22main%22%7D
			// return c.OriginalURL() //这里不用是因为本地开发中有问题
			return string(c.Request().Header.RequestURI()) //默认值是c.Path(),同一个接口，不同请求参数默认也会缓存,所以这里改成c.OriginalURL()
		},
	})
}

// NewCacheControl 创建前端强缓存的中间件
func NewCacheControl(maxAge int) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderCacheControl, fmt.Sprintf("public, max-age=%d", maxAge))
		return c.Next()
	}
}
