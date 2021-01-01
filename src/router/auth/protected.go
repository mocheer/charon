package auth

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/mocheer/charon/src/core/res"
)

// Protected protect routes
func Protected(config jwtware.Config) fiber.Handler {
	return jwtware.New(config)
}

// GlobalProtected 全局的认证handler
var GlobalProtected = jwtware.New(jwtware.Config{
	SigningKey:   SigningKey,
	ErrorHandler: jwtError,
})

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return res.ResultError(c, fiber.StatusBadRequest, "Missing or malformed JWT", nil)
	}
	return res.ResultError(c, fiber.StatusUnauthorized, "Invalid or expired JWT", nil)
}
