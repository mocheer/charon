package reptile

import (
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/res"
	"github.com/mocheer/vesta"
)

func Use(api fiber.Router) {
	router := api.Group("/reptile")
	//
	router.Get("/", reptileHandle)
	router.Post("/", reptileHandle)
	// 网页截图
	router.Get("/image", repliteScreen)
}

// 抓取数据
func reptileHandle(c *fiber.Ctx) error {
	args := &reptileArgs{}
	if err := c.QueryParser(args); err != nil {
		return err
	}
	data := vesta.New().Nav(args.Url).GetValue(args.Script)
	return res.JSON(c, data)
}

// 抓取图片
func repliteScreen(c *fiber.Ctx) error {
	args := &reptileArgs{}
	if err := c.QueryParser(args); err != nil {
		return err
	}
	v := vesta.New().Nav(args.Url)
	// W 默认值800
	// H 默认值600
	v.Viewport(args.W, args.H)
	if args.WaitQuery != "" {
		v.WaitQuery(args.WaitQuery)
	}
	v.Wait(300 * time.Millisecond)
	data := v.GetScreen()
	//
	switch args.F {
	case "base64":
		return c.SendString(base64.StdEncoding.EncodeToString(data))
	default:
		return res.PNG(c, data)
	}
}
