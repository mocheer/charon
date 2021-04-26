package auth

import (
	"github.com/mocheer/charon/src/global"
	"github.com/mocheer/pluto/fn"
)

// SigningKey 密钥
var SigningKey = func() []byte {
	return fn.StringBytes(global.Config.Name)
}
