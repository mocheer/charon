package auth

import (
	"github.com/mocheer/charon/cts"
	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/orm/tables"
	"github.com/mocheer/pluto/ecc"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/bcrypt"
)

// DecodeCliper 解析前端密文
func DecodeCliper(data string) gjson.Result {
	plainText := ecc.RSA_DecodeJSEncrypt(data, cts.RSA_PrivatePemPath)
	return gjson.Parse(plainText)
}

// CheckPasswordHash 对比hash密码和输入的密码
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// getUserByUsername 根据用户名获取用户信息
func getUserByUsername(userName string) (*tables.User, error) {
	var user tables.User
	if err := global.DB.Where(&tables.User{Name: userName}).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// hashAndSalt 加盐加密
func hashAndSalt(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash)
}
