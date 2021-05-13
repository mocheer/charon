package req

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

//
func Parse(c *fiber.Ctx, out interface{}) {
	if err := c.BodyParser(out); err != nil {
		panic(fmt.Sprintf("参数有误：%s", err.Error()))
	}
}
