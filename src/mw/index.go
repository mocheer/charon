package mw

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/mocheer/charon/src/global"
	"github.com/mocheer/pluto/clock"
	"github.com/mocheer/pluto/fn"
)

// initMW 初始化middleware中间件
func initMW() {
	//
	SigningKey = fn.StringBytes(global.Config.Name)
	// token路由守卫
	GlobalProtected = jwtware.New(jwtware.Config{
		SigningKey:   SigningKey,
		ErrorHandler: jwtError,
	})
}

// Use 使用所有中间件
func Use(app *fiber.App) {
	// 初始化
	initMW()
	// /debug/pprof/
	app.Use(pprof.New())
	// 日志中间件
	app.Use(logger.New(logger.Config{
		Output: os.Stdout,
	}))
	// 协商缓存
	app.Use(etag.New())
	// 插件有使用顺序，且顺序非常重要，比如说cache需要放到compress后面(这个在2.2.4之后版本已支持)，compresss需要放到业务路由前面等
	// recover 中间件，防止因为某个路由的错误导致整个应用崩溃
	// 发生错误时状态码为500，而且会将错误数据返回到前端
	app.Use(recover.New())
	// cors
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0, //缓存，单位秒
	}))

	//
	app.Use(func(c *fiber.Ctx) error {
		if strings.HasSuffix(c.Path(), ".js.map") {
			return c.SendString(`{
			"version": 3,
			"sources": [],
			"names": [],
			"mappings": "",
			"file": "",
			"sourcesContent": [],
			"sourceRoot": ""
		}`)
		}
		return c.Next()
	})
	//
	for name, config := range global.Config.Static {
		app.Static(name, config.Dir, fiber.Static{
			Compress:  true,        //
			ByteRange: true,        //
			Browse:    false,       // 是否访问目录时列出文件列表
			MaxAge:    clock.Month, // 缓存时间，单位秒
			Index:     "index.html",
		})
	}

	// 重定向
	app.Get("/v/*", func(c *fiber.Ctx) error {
		// maxAge := strconv.FormatUint(86400*30, 10)
		// c.Set(fiber.HeaderCacheControl, "public, max-age="+maxAge)
		return c.SendFile("./public/index.html")
	})

	// 重定向
	protected := jwtware.New(jwtware.Config{
		SigningKey: SigningKey,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			cookie := new(fiber.Cookie)
			cookie.Name = "c_url"
			cookie.Value = string(c.Request().Header.RequestURI())
			cookie.Expires = time.Now().Add(5 * time.Second)
			c.Cookie(cookie)
			return c.Redirect("/v/studio/login")
		},
		TokenLookup: "cookie:t",
	})
	//
	app.Get("/web/*", protected, func(c *fiber.Ctx) error {
		return c.SendFile("./web/index.html", false)
	})

	// 不支持压缩中间件，所以只能放到这个中间件前面实例化
	//app.Use("/docs", swagger.Handler) // default

	// 强缓存=>不是所有的请求都需要强缓存
	// app.Use(store.GlobalCache)
	//
	// 压缩中间件
	// 为什么 localhost 和127.0.0.1的请求时是br，而192.168.117.215或者其他远程服务器是gzip
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

}
