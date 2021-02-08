package auth

import (
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/src/core/res"
	"github.com/mocheer/charon/src/global"
	"github.com/mocheer/charon/src/models/tables"
)

// Use 初始化 auth 路由
// @see https://github.com/gofiber/recipes/tree/master/auth-jwt
func Use(api fiber.Router) {
	router := api.Group("/auth")
	router.Post("/login", login)
	router.Post("/signup", signup)
	router.Get("/info", GlobalProtected, restricted)
}

// login
func login(c *fiber.Ctx) error {

	var input LoginInput

	if err := c.BodyParser(&input); err != nil {
		return res.ResultError(c, fiber.StatusBadRequest, "Error on login request", err)
	}
	username := input.Username
	password := input.Password

	user, err := getUserByUsername(username)
	if err != nil {
		return res.ResultError(c, fiber.StatusUnauthorized, "Error on username", err)
	}

	// 这里直接判断原始密码有问题
	if !CheckPasswordHash(password, user.Password) {
		return res.ResultError(c, fiber.StatusUnauthorized, "Invalid password", nil)
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString(SigningKey)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	// 用于一些服务的路由验证
	cookie := new(fiber.Cookie)
	cookie.Name = "t"
	cookie.Value = t
	cookie.Expires = time.Now().Add(72 * time.Hour)
	// Set cookie
	c.Cookie(cookie)

	return res.ResultOK(c, t)
}

// signup 注册
func signup(c *fiber.Ctx) error {
	var input LoginInput

	if err := c.BodyParser(&input); err != nil {
		return res.ResultError(c, fiber.StatusBadRequest, "参数有误", err)
	}
	password := hashAndSalt(input.Password)
	query := global.Db.Create(tables.User{Name: input.Username, Password: password})
	//
	if query.Error != nil {
		return res.ResultError(c, fiber.StatusInternalServerError, "注册失败", query.Error)
	}

	return res.ResultOK(c, true)
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
