package global

import "github.com/mocheer/charon/src/core/fs"

type staticConfig struct {
	// 前端资源路由名称
	Name string
	// 前端资源路由模式
	Mode string
	// 静态资源目录
	Dir string
}

// appConfig 应用配置
type appConfig struct {
	// 名称
	Name string
	// 应用服务信息
	DisplayName string
	// 应用模式：dev|production
	Mode string
	// 监听端口
	Port string
	// 静态资源服务
	Static map[string]staticConfig
	// 数据库连接串
	DbDSN string
}

// ReadJSON 从json文件中读取配置数据
func (e *appConfig) ReadJSON(path string) error {
	return fs.ReadJSON(path, e)
}
