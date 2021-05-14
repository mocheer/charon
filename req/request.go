package req

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/orm"
)

func Engine() *orm.Wrapper {
	ctx := global.DB
	return &orm.Wrapper{
		Ctx: ctx,
	}
}

//
func MustParseSelectArgs(c *fiber.Ctx) *orm.SelectArgs {
	var args = &orm.SelectArgs{}
	if err := c.QueryParser(args); err != nil {
		panic(fmt.Sprintf("参数有误：%s", err.Error()))
	}
	return args
}
