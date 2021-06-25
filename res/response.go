package res

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/global"
	"github.com/mocheer/pluto/ts/img"
)

// Response 返回的数据结构
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  interface{} `json:"msg"`
}

// Result 约定返回的数据格式
func Result(c *fiber.Ctx, code int, data interface{}, msg string) error {
	return c.Status(code).JSON(&Response{
		code,
		data,
		msg,
	})
}

// JSON 返回成功
func JSON(c *fiber.Ctx, data interface{}) error {
	return Result(c, fiber.StatusOK, data, "")
}

// Error 返回错误信息
func Error(c *fiber.Ctx, data string, err error) error {
	if global.IS_DEV {
		fmt.Println(err)
	}
	return Result(c, fiber.StatusInternalServerError, data, err.Error())
}

// Warn 警告
func Warn(c *fiber.Ctx, data string) error {
	return Result(c, fiber.StatusInternalServerError, data, "")
}

// Image 返回图片
func Image(c *fiber.Ctx, p *img.Img) error {
	bs, err := p.ToBytes()
	if err != nil {
		return err
	}
	return c.Type(p.Type).Send(bs)
}

// PNG 返回图片信息
func PNG(c *fiber.Ctx, data []byte) error {
	return c.Type(img.PNG).Send(data)
}
