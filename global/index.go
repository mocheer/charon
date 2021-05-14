package global

import (
	"github.com/mocheer/charon/cts"
	"github.com/mocheer/pluto/fs"
)

// 开发模式
var IS_DEV bool

func Init(mode string) {
	IS_DEV = mode == "" || mode == cts.DevMode
	// 初始化日志
	Log.Init()
	// 初始化配置
	Config = &AppConfig{
		Name: "charon",
		Port: ":9212",
		Static: []StaticConfig{
			{
				Name: "/web",
				Mode: "history",
				Dir:  "./assets",
			},
		},
	}
	// 读取应用配置
	err := fs.ReadJSON(cts.AppConfigPath, Config)
	if err != nil {
		Log.Info("配置文件不存在，重新生成")
		err = fs.SaveJSON(cts.AppConfigPath, Config)
		if err != nil {
			Log.Info("生成配置文件失败", err)
		}
	}
	// 连接数据库
	DB, err = openDB()
	if err != nil {
		Log.Error("数据库连接失败", err)
	}

	initRSA()
	//
	if IS_DEV {
		Log.Info("启动开发模式成功")
	} else {
		Log.Info("启动产品模式成功")
	}
}
