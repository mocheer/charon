package router

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/mocheer/charon/src/global"
	"github.com/mocheer/pluto/fs"

	"github.com/mocheer/charon/src/router/arcgis"
	"github.com/mocheer/charon/src/router/auth"
	"github.com/mocheer/charon/src/router/dmap"
	"github.com/mocheer/charon/src/router/model"
	"github.com/mocheer/charon/src/router/pipal"
	"github.com/mocheer/charon/src/router/proxies"
	"github.com/mocheer/charon/src/router/query"
	"github.com/mocheer/charon/src/router/upload"

	jwtware "github.com/gofiber/jwt/v2"
)

// useRouter
func useRouter(app *fiber.App) {
	api := app.Group("/api") // "/api" 这里不用api是为了与ec等框架做区分 => "/con"
	v1 := api.Group("/v1")
	//
	pipal.Use(v1)
	query.Use(v1)
	upload.Use(v1)
	model.Use(v1)
	auth.Use(v1)
	proxies.Use(v1)
	//
	arcgis.Use(v1)
	dmap.Use(v1)
}

// Init 初始化路由
func Init() {
	app := fiber.New()
	// /debug/pprof/
	app.Use(pprof.New())
	// 日志中间件
	app.Use(logger.New(logger.Config{
		Output: os.Stdout,
	}))
	// 插件有使用顺序，且顺序非常重要，比如说cache需要放到compress后面(这个在2.2.4之后版本已支持)，compresss需要放到业务路由前面等
	// recover 中间件，防止因为某个路由的错误导致整个应用崩溃
	// 发生错误时状态码为500，而且会将错误数据返回到前端
	app.Use(recover.New())
	// 图标
	if fs.IsExist("./public/favicon.ico") {
		app.Use(favicon.New(favicon.Config{
			File: "./public/favicon.ico",
		}))
	}
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

	for name, config := range global.Config.Static {
		app.Static(name, config.Dir, fiber.Static{
			Compress:  true,       //
			ByteRange: true,       //
			Browse:    false,      // 是否访问目录时列出文件列表
			MaxAge:    86400 * 30, // 缓存时间，单位秒，86400s = 1day
			Index:     "index.html",
		})
	}

	// filepathNames, err := filepath.Glob(filepath.Join("./public", "*"))
	// if err == nil {
	// 	for i := range filepathNames {
	// 		fmt.Println(filepathNames[i]) //打印path
	// 	}
	// }

	// 不支持压缩中间件，所以只能放到这个中间件前面实例化
	//app.Use("/docs", swagger.Handler) // default

	// 全局缓存
	// app.Use(store.GlobalCache)
	//
	// 压缩中间件
	// 为什么 localhost 和127.0.0.1的请求时是br，而192.168.117.215或者其他远程服务器是gzip
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	//
	useRouter(app)

	// 重定向
	app.Get("/v/*", func(c *fiber.Ctx) error {
		// maxAge := strconv.FormatUint(86400*30, 10)
		// c.Set(fiber.HeaderCacheControl, "public, max-age="+maxAge)
		return c.SendFile("./public/index.html", false)
	})
	// 重定向
	protected := jwtware.New(jwtware.Config{
		SigningKey: auth.SigningKey,
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
	// 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("not found")
	})
	//
	log.Fatal(app.Listen(global.Config.Port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	// app.ListenTLS(":443", "./cert.pem", "./cert.key");//2.3.0
}
