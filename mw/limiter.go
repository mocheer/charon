package mw

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// NewLimiter 限制接口访问，用于节流
func NewLimiter(conf limiter.Config) func(*fiber.Ctx) error {
	return limiter.New(conf)
}
