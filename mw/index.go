package mw

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/mocheer/charon/global"
	"github.com/mocheer/pluto/fn"
	"github.com/mocheer/pluto/ts/clock"
)

// Use 使用所有中间件
func Use(app *fiber.App) {
	SigningKey = fn.String2Bytes(global.Config.Name)
	// 日志中间件
	app.Use(logger.New(logger.Config{
		Output: os.Stdout,
	}))
	// 已经有token验证了
	// app.Use(csrf.New(csrf.Config{
	// 	Next: func(c *fiber.Ctx) bool {
	// 		return c.Path() == `/api/v1/auth/login`
	// 	},
	// }))
	// 安全中间件，包含xss、xframe、contenttype等方面的漏洞防御
	app.Use(helmet.New())

	// 插件有使用顺序，且顺序非常重要，比如说cache需要放到compress后面(这个在2.2.4之后版本已支持)，compresss需要放到业务路由前面等
	// recover 中间件，防止因为某个路由的错误导致整个应用崩溃
	// 发生错误时状态码为500，而且会将错误数据返回到前端
	app.Use(recover.New(recover.Config{
		EnableStackTrace: global.IS_DEV,
	}))
	// cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0, //缓存，单位秒
	}))
	//
	if !global.IS_DEV {
		// 生产环境下拦截.js.map文件
		app.Use(HideJSMap)
		// 协商缓存
		app.Use(etag.New())
	} else {
		// 开发环境下支持pprof调试
		app.Use(pprof.New())
	}

	for index, config := range global.Config.Static {

		app.Static(config.Name, config.Dir, fiber.Static{
			Compress:  true,       //
			ByteRange: true,       //
			Browse:    false,      // 是否访问目录时列出文件列表
			MaxAge:    clock.Week, // 强缓存时间，单位秒
			Index:     "index.html",
		})
		//
		indexHTML := filepath.Join(config.Dir, "index.html")
		routeName := fmt.Sprintf("%s/*", config.Name)
		if index == 0 {
			app.Get(routeName, func(c *fiber.Ctx) error {
				return c.SendFile(indexHTML)
			})
		} else {
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
			app.Get(routeName, protected, func(c *fiber.Ctx) error {
				return c.SendFile(indexHTML)
			})
		}
	}

	// 压缩中间件
	// 为什么 localhost 和 127.0.0.1 的请求时是br，而 192.168.117.215 或者其他远程服务器是 gzip
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

}

// UseProtected
func UseProtected(r fiber.Router) {
	// jwt token认证守卫
	r.Use(jwtware.New(jwtware.Config{
		Filter: func(c *fiber.Ctx) bool {
			return c.Method() == "GET" || strings.HasPrefix(c.Path(), `/api/v1/query/raw`)
		},
		SigningKey:   SigningKey,
		ErrorHandler: jwtError,
	}))
}
