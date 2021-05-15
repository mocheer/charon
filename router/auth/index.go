package auth

import (
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/mocheer/charon/cts"
	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/mw"
	"github.com/mocheer/charon/orm/tables"
	"github.com/mocheer/charon/res"
)

// Use 初始化 auth 路由
// @see https://github.com/gofiber/recipes/tree/master/auth-jwt
func Use(api fiber.Router) {
	router := api.Group("/auth")
	// 登录
	router.Post("/login", mw.NewLimiter(limiter.Config{Max: 3, Expiration: 10 * time.Second}), login)
	// 注册
	router.Post("/signup", signup)
	// 获取用户信息
	router.Get("/info", mw.Protected(), getUserInfo)
	// 获取用于加密用的rsa公钥，用于前端加密
	router.Get("/rsa", res.HandleTextFile(cts.RSA_PublicPemPath))
}

// login 登录
func login(c *fiber.Ctx) error {
	var input LoginInput
	//
	if err := c.BodyParser(&input); err != nil {
		return res.Result(c, fiber.StatusBadRequest, "Error on login request", err)
	}
	//
	plain := DecodeCliper(input.Cipher)
	username := plain.Get("username").String()
	password := plain.Get("password").String()
	//
	user, err := getUserByUsername(username)
	if err != nil {
		return res.Result(c, fiber.StatusUnauthorized, "Error on username", err)
	}

	// 这里直接判断原始密码有问题
	if !CheckPasswordHash(password, user.Password) {
		return res.Result(c, fiber.StatusUnauthorized, "Invalid password", nil)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	//
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	//
	t, err := token.SignedString(mw.SigningKey)
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

	return res.JSON(c, t)
}

// signup 注册
func signup(c *fiber.Ctx) error {
	var input LoginInput
	//
	if err := c.BodyParser(&input); err != nil {
		return res.Result(c, fiber.StatusBadRequest, "参数有误", err)
	}

	plain := DecodeCliper(input.Cipher)
	username := plain.Get("username").String()
	password := plain.Get("password").String()
	//
	if len(password) < 6 {
		return res.Result(c, fiber.StatusBadRequest, "参数有误：密码不应小于6位", nil)
	}
	//
	query := global.DB.Create(tables.User{Name: username, Password: hashAndSalt(password)})
	//
	if query.Error != nil {
		return res.Result(c, fiber.StatusInternalServerError, "注册失败", query.Error)
	}
	//
	return res.JSON(c, true)
}

// getUserInfo 获取token中的用户信息
func getUserInfo(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)

	claims := user.Claims.(jwt.MapClaims)
	return c.JSON(claims)
}
