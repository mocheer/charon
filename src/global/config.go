package global

import "github.com/gofiber/fiber/v2/middleware/cors"

// StaticConfig 静态资源配置
type StaticConfig struct {
	// 名称，暂定，可根据名称，查找到静态资源的相关信息
	Name string
	// 路由前缀
	Prefix string
	// 路由模式 history | hash
	Mode string
	// 静态资源根目录
	Dir string
}

// AppConfig 应用配置
type AppConfig struct {
	// 名称
	Name string
	// 应用服务信息
	DisplayName string
	// 应用模式：dev|production
	Mode string
	// 监听端口
	Port string
	// 静态资源服务
	Static map[string]StaticConfig
	//
	Cors cors.Config
	// DSN(Data Source Name)数据源：数据库连接串
	DSN string
}
