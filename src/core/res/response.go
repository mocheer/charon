package res

import (
	"github.com/gofiber/fiber/v2"
)

// Response 返回的数据结构
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  interface{} `json:"msg"`
}

// Result 约定返回的数据格式
func Result(c *fiber.Ctx, code int, data interface{}, msg interface{}) error {
	return c.Status(code).JSON(Response{
		code,
		data,
		msg,
	})
}

// ResultOK 返回成功
func ResultOK(c *fiber.Ctx, data interface{}) error {
	return Result(c, fiber.StatusOK, data, "")
}

// ResultError 返回错误信息
func ResultError(c *fiber.Ctx, data string, err error) error {
	return Result(c, fiber.StatusInternalServerError, data, err)
}

func ResultImage(c *fiber.Ctx, data interface{}) error {
	return nil
}
