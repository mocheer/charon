package mw

import (
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/mocheer/charon/res"
)

// SigningKey 密钥
var SigningKey []byte

//
var protector func(*fiber.Ctx) error

// Protected protect routes
func Protected() fiber.Handler {
	if protector == nil {
		protector = jwtware.New(jwtware.Config{
			SigningKey:   SigningKey,
			ErrorHandler: jwtError})
	}
	return protector
}

// PermissProtectd 特殊权限认证 role=1 有权限
func PermissProtectd(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	role := claims["role"].(float64)
	if role == 1 {
		return c.Next()
	}
	return res.Result(c, fiber.StatusBadRequest, "没有权限", "")
}

// jwtError jwt 认证失败
func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return res.Result(c, fiber.StatusBadRequest, err.Error(), "")
	}
	return res.Result(c, fiber.StatusUnauthorized, "Invalid or expired JWT", "")
}
