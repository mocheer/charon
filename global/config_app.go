package global

// AppConfig 应用配置
type AppConfig struct {
	// 名称
	Name string
	// 应用服务信息
	DisplayName string
	// 是否启用https（前提条件：有申请对应的ssl证书）
	IsHTTPS bool
	// 监听端口
	Port string
	// 静态资源服务
	Static []StaticConfig
	// DSN(Data Source Name)数据源：数据库连接串
	DSN string
}

// Config 全局配置
var Config *AppConfig

// FirstStaticDir 第一个静态资源目录
func (c *AppConfig) FirstStaticDir() string {
	return c.Static[0].Dir
}
