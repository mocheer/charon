package res

import (
	"github.com/gofiber/fiber/v2"
)

// Response 返回的数据结构
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// Result 约定返回的数据格式
func Result(c *fiber.Ctx, data interface{}, code int, msg string) error {
	// 开始时间
	return c.JSON(Response{
		code,
		data,
		msg,
	})
}

// ResultOK 返回成功
func ResultOK(c *fiber.Ctx, data interface{}) error {
	// 开始时间
	return Result(c, data, fiber.StatusOK, "")
}

// ResultError 返回错误信息
func ResultError(c *fiber.Ctx, code int, msg string, err error) error {
	return Result(c.Status(code), err, code, msg)
}
