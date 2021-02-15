package auth

import (
	"github.com/mocheer/charon/src/global"
	"github.com/mocheer/charon/src/models/tables"
	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordHash 对比hash密码和输入的密码
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// getUserByUsername 根据用户名获取用户信息
func getUserByUsername(userName string) (*tables.User, error) {
	var user tables.User
	if err := global.Db.Where(&tables.User{Name: userName}).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func hashAndSalt(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash)
}