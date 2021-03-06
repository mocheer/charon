package req

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/orm"
)

func Engine() *orm.Wrapper {
	ctx := global.DB
	if global.IS_DEV {
		ctx = ctx.Debug() // debug 会启动一个新的会话，不能在model之后
	}
	return &orm.Wrapper{
		Ctx: ctx,
	}
}

//
func MustParseArgs(c *fiber.Ctx, args interface{}) interface{} {
	if err := c.QueryParser(args); err != nil {
		panic(fmt.Sprintf("URL参数有误：%s", err.Error()))
	}
	if c.Method() == "POST" {
		if err := c.BodyParser(args); err != nil {
			panic(fmt.Sprintf("POST参数有误：%s", err.Error()))
		}
	}
	return args
}
