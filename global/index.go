package global

import (
	"github.com/mocheer/pluto/fs"
)

func Init() {
	// 初始化日志
	Log.Init()
	// 初始化配置
	Config = &AppConfig{
		Name: "charon",
		Mode: "production",
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
	err := fs.ReadJSON(AppConfigPath, Config)
	if err != nil {
		Log.Info("配置文件不存在，重新生成")
		err = fs.SaveJSON(AppConfigPath, Config)
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
	Log.Info("启动成功")
}
