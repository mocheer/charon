package auth

import (
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/mocheer/charon/src/res"
)

// Protected protect routes
func Protected(config jwtware.Config) fiber.Handler {
	return jwtware.New(config)
}

// GlobalProtected 全局的认证handler
var GlobalProtected = jwtware.New(jwtware.Config{
	SigningKey:   SigningKey(),
	ErrorHandler: jwtError,
})

// PermissProtectd 特殊权限认证
func PermissProtectd(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	role := claims["role"].(float64)
	if role == 1 {
		return c.Next()
	}
	return res.Result(c, fiber.StatusBadRequest, "没有权限", nil)
}

// jwtError jwt 认证失败
func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return res.Result(c, fiber.StatusBadRequest, "Missing or malformed JWT", nil)
	}
	return res.Result(c, fiber.StatusUnauthorized, "Invalid or expired JWT", nil)
}
