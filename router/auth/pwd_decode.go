package auth

import (
	"encoding/base64"

	"github.com/mocheer/pluto/fn"
)

// DecodePassword 解析加密后的用户密码  btoa(password)).split("").map(e=>String.fromCharCode(e.charCodeAt(0)+1)).reverse().join("")
func DecodePassword(password string) (string, error) {
	str := fn.ReverseStringOffset(password, -1)
	strByte, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return fn.BytesString(strByte), nil
}
