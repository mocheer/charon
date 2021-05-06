package mw

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

//
func NewLimiter(conf limiter.Config) func(*fiber.Ctx) error {
	return limiter.New(conf)
}
