package global

import (
	"github.com/mocheer/charon/cts"
	"github.com/mocheer/pluto/ecc"
	"github.com/mocheer/pluto/fs"
)

// initRSA 初始化rsa生成密钥文件
func initRSA() {
	if !fs.IsExist(cts.RSA_PrivatePemPath) {
		err := ecc.RSA_GenPemFiles(cts.RSA_Dir, 2048)
		if err == nil {
			Log.Info("成功生成RSA密钥文件")
		} else {
			Log.Warn("生成RSA密钥文件失败")
		}
	}
}
