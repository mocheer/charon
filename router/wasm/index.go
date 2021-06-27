package wasm

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/mw"
)

// {GoRoot}/misc/wasm/wasm_exec.js
//go:embed wasm_exec.js
var wasmExec_js []byte

func Use(api fiber.Router) {
	router := api.Group("/wasm")
	router.Get("/", mw.Cache, func(c *fiber.Ctx) error {
		return c.Send(wasmExec_js)
	})
}
